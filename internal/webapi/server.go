package webapi

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/internal/authentication"
	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/tunnel"
	"github.com/gpuctl/gpuctl/internal/types"
)

type Server struct {
	mux *femto.Femto
	api *Api
}
type Api struct {
	DB         database.Database
	tunnelConf tunnel.Config
}

type APIAuthCredientals struct {
	Username string
	Password string
}

func NewServer(db database.Database, auth authentication.Authenticator[APIAuthCredientals], tunnelConf tunnel.Config) *Server {
	mux := new(femto.Femto)
	api := &Api{db, tunnelConf}

	femto.OnGet(mux, "/api/stats/all", api.AllStatistics)
	femto.OnGet(mux, "/api/stats/offline", api.HandleOfflineMachineRequest)
	femto.OnGet(mux, "/api/stats/since_last_seen", api.durationDelta)

	// Set up authentication and logging-out endpoint
	femto.OnPost(mux, "/api/admin/auth", func(packet APIAuthCredientals, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return api.Authenticate(auth, packet, r, l)
	})
	femto.OnGet(mux, "/api/admin/logout", func(r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return api.LogOut(auth, r, l)
	})

	// Authenticated API endpoints
	femto.OnPost(mux, "/api/admin/add_workstation", authentication.AuthWrapPost(auth, api.addMachine))
	femto.OnPost(mux, "/api/admin/stats/modify", authentication.AuthWrapPost(auth, api.modifyMachineInfo))
	femto.OnPost(mux, "/api/admin/rm_workstation", authentication.AuthWrapPost(auth, api.removeMachine))
	//femto.OnPost(mux, "/api/admin/attach_file", authentication.AuthWrapPost(auth, api.attachFile))
	femto.OnPost(mux, "/api/admin/attach_file", api.attachFile)
	femto.OnGet(mux, "/api/admin/confirm", authentication.AuthWrapGet(auth, func(r *http.Request, l *slog.Logger) (*femto.Response[UsernameReminder], error) {
		return api.ConfirmAdmin(auth, r, l)
	}))

	return &Server{mux, api}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// TODO: Maybe unset in Caddyfile???
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	s.mux.ServeHTTP(w, r)
}

// This function involves a lot of weird unwrapping
// TODO: See if we can get the database layer to do it for us
func (a *Api) AllStatistics(r *http.Request, l *slog.Logger) (*femto.Response[broadcast.Workstations], error) {
	data, err := a.DB.LatestData()

	if err != nil {
		return nil, err
	}
	if data == nil {
		// dont just return nil, which would not be marshalled properly
		return &femto.Response[broadcast.Workstations]{Body: broadcast.Workstations{}}, nil
	}

	return femto.Ok(data)
}

func (a *Api) Authenticate(auth authentication.Authenticator[APIAuthCredientals], packet APIAuthCredientals, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	// Check if credientals are correct
	token, err := auth.CreateToken(packet)

	if errors.Is(err, authentication.InvalidCredentialsError) || errors.Is(err, authentication.NotAuthenticatedError) {
		return &femto.Response[types.Unit]{Status: http.StatusUnauthorized}, nil
	}

	if err != nil {
		return nil, err
	}

	cookies := []http.Cookie{{
		Name:     authentication.TokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}}

	return &femto.EmptyBodyResponse{Cookies: cookies, Status: http.StatusAccepted}, nil
}

func (a *Api) LogOut(auth authentication.Authenticator[APIAuthCredientals], r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	token, err := r.Cookie(authentication.TokenCookieName)
	if err != nil {
		return nil, err
	}
	auth.RevokeToken(token.Value)
	return femto.Ok(types.Unit{})
}

// TODO
func (a *Api) addMachine(machine broadcast.NewMachine, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	l.Info("Tried to create machine", "host", machine.Hostname, "group", machine.Group)

	_, err := a.onboard(broadcast.OnboardReq{Hostname: machine.Hostname}, r, l)
	if err != nil {
		return nil, err
	}
	modify := broadcast.ModifyMachine{Hostname: machine.Hostname, Group: machine.Group}
	_, err = a.modifyMachineInfo(modify, r, l)
	if err != nil {
		return nil, err
	}

	return femto.Ok(types.Unit{})
}

func (a *Api) attachFile(attach broadcast.AttachFile, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	err := a.DB.AttachFile(attach)
	if err != nil {
		return nil, err
	}
	return femto.Ok(types.Unit{})
}

func (a *Api) removeMachine(rm broadcast.RemoveMachineInfo, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	err := a.deboard(rm, r, l)
	if err != nil {
		return nil, err
	}
	err = a.DB.RemoveMachine(broadcast.RemoveMachine{Hostname: rm.Hostname})
	if err != nil {
		return nil, err
	}
	return femto.Ok(types.Unit{})
}

type UsernameReminder struct {
	Username string `json:"username"`
}

func (a *Api) ConfirmAdmin(auth authentication.Authenticator[APIAuthCredientals], r *http.Request, l *slog.Logger) (*femto.Response[UsernameReminder], error) {
	c, err := r.Cookie(authentication.TokenCookieName)
	if err != nil {
		return &femto.Response[UsernameReminder]{Status: http.StatusUnauthorized}, nil
	}
	u, err := auth.CheckToken(c.Value)
	if err != nil {
		return &femto.Response[UsernameReminder]{Status: http.StatusUnauthorized}, nil
	}
	return femto.Ok(UsernameReminder{Username: u})
}

func (a *Api) durationDelta(r *http.Request, l *slog.Logger) (*femto.Response[[]broadcast.DurationDeltas], error) {
	const nanosInSecond = 1e9

	latest, err := a.DB.LastSeen()

	if err != nil {
		return nil, err
	}

	var deltas []broadcast.DurationDeltas

	now := time.Now().Unix() / nanosInSecond

	for idx := range latest {
		then_s := latest[idx].LastSeen / nanosInSecond

		deltas = append(deltas, broadcast.DurationDeltas{
			Hostname: latest[idx].Hostname,
			Delta:    now - then_s,
		})
	}

	return femto.Ok(deltas)
}

func (a *Api) modifyMachineInfo(info broadcast.ModifyMachine, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	l.Info("Tried to modify machine", "host", info.Hostname, "changes", info)

	err := a.DB.UpdateMachine(info)
	if err != nil {
		return nil, err
	}

	return femto.Ok(types.Unit{})
}
