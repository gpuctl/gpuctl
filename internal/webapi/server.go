package webapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/internal/authentication"
	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/types"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	mux *femto.Femto
	api *Api
}
type Api struct {
	DB          database.Database
	onboardConf OnboardConf
}

type APIAuthCredientals struct {
	Username string
	Password string
}

type OnboardConf struct {
	// The login to run the satellite on other machines as
	Username string
	// The directory to store the satellite binary on remotes as
	DataDir string
	// The configuration to install on the remote.
	RemoteConf config.SatelliteConfiguration

	// SSH Options.
	Signer      ssh.Signer
	KeyCallback ssh.HostKeyCallback
}

func NewServer(db database.Database, auth authentication.Authenticator[APIAuthCredientals], onboard OnboardConf) *Server {
	mux := new(femto.Femto)
	api := &Api{db, onboard}

	femto.OnGet(mux, "/api/stats/all", api.AllStatistics)
	femto.OnGet(mux, "/api/stats/offline", api.HandleOfflineMachineRequest)
	femto.OnGet(mux, "/api/stats/since_last_seen", api.durationDelta)

	// Set up authentication endpoint
	femto.OnPost(mux, "/api/admin/auth", func(packet APIAuthCredientals, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return api.Authenticate(auth, packet, r, l)
	})

	// Authenticated API endpoints
	femto.OnPost(mux, "/api/admin/add_workstation", api.addMachine)
	femto.OnPost(mux, "/api/admin/stats/modify", api.modifyMachineInfo)
	femto.OnPost(mux, "/api/admin/rm_workstation", api.removeMachine)
	femto.OnGet(mux, "/api/admin/confirm", func(r *http.Request, l *slog.Logger) (*femto.Response[UsernameReminder], error) {
		return api.confirmAdmin(auth, r, l)
	})

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

	var ws []broadcast.WorkstationData
	for _, machine := range data {
		if len(machine.Stats) == 0 {
			continue
		}

		gpus := make([]broadcast.OldGPUStatSample, 0)
		for i := range machine.Stats {
			gpus = append(gpus, ZipStats(
				machine.Hostname,
				machine.GPUInfos[i],
				machine.Stats[i],
			))
		}

		ws = append(ws, broadcast.WorkstationData{
			Name: machine.Hostname,
			Gpus: gpus,
		})
	}

	result := broadcast.Workstations{{Name: "Shared", WorkStations: ws}}
	return femto.Ok(result)
}

func (a *Api) Authenticate(auth authentication.Authenticator[APIAuthCredientals], packet APIAuthCredientals, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	// Check if credientals are correct
	token, err := auth.CreateToken(packet)

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

func (a *Api) confirmAdmin(auth authentication.Authenticator[APIAuthCredientals], r *http.Request, l *slog.Logger) (*femto.Response[UsernameReminder], error) {
	c, err := r.Cookie(authentication.TokenCookieName)
	slog.Info("TEST", "Cookie", c, "Err", err)
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

// Bodge together stats and contextual data to make OldGpuStats
func ZipStats(host string, info uplink.GPUInfo, stat uplink.GPUStatSample) broadcast.OldGPUStatSample {
	return broadcast.OldGPUStatSample{
		Hostname: host,
		// info from GPUInfo
		Uuid:          info.Uuid,
		Name:          info.Name,
		Brand:         info.Brand,
		DriverVersion: info.DriverVersion,
		MemoryTotal:   info.MemoryTotal,
		// info from GPUStatSample
		MemoryUtilisation: stat.MemoryUtilisation,
		GPUUtilisation:    stat.GPUUtilisation,
		MemoryUsed:        stat.MemoryUsed,
		FanSpeed:          stat.FanSpeed,
		Temp:              stat.Temp,
		MemoryTemp:        stat.MemoryTemp,
		GraphicsVoltage:   stat.GraphicsVoltage,
		PowerDraw:         stat.PowerDraw,
		GraphicsClock:     stat.GraphicsClock,
		MaxGraphicsClock:  stat.MaxGraphicsClock,
		MemoryClock:       stat.MemoryClock,
		MaxMemoryClock:    stat.MaxMemoryClock,
	}
}
