package webapi

import (
	"fmt"
	"log/slog"
	"net/http"

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

func makeAuthCookie(token string) string {
	return fmt.Sprintf("token=%s; Path=/; HttpOnly; Secure; SameSite=Strict", token)
}

// ! Key change
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

	// Authenticated API endpoints

	femto.OnPost(mux, "/api/admin/auth", func(packet APIAuthCredientals, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
		return api.Authenticate(auth, packet, r, l)
	})
	femto.OnPost(mux, "/api/admin/add_workstation", authentication.AuthWrapPost(auth, api.newMachine))
	femto.OnPost(mux, "/api/machines/addinfo", authentication.AuthWrapPost(auth, api.addInfo))
	femto.OnPost(mux, "/api/onboard", authentication.AuthWrapPost(auth, api.onboard))

	return &Server{mux, api}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
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

	result := []broadcast.WorkstationGroup{{Name: "Shared", WorkStations: ws}}
	response := femto.Ok[broadcast.Workstations](result)
	return &response, nil
}

func (a *Api) Authenticate(auth authentication.Authenticator[APIAuthCredientals], packet APIAuthCredientals, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	// Check if credientals are correct
	token, err := auth.CreateToken(packet)

	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	headers["Set-Cookie"] = makeAuthCookie(token)
	return &femto.EmptyBodyResponse{Headers: headers, Status: http.StatusAccepted}, nil
}

// Create a new machine
func (a *Api) newMachine(machine broadcast.NewMachine, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	l.Info("Tried to create machine", "host", machine.Hostname, "group", machine.Group)
	response := femto.Ok[types.Unit](types.Unit{})
	return &response, a.DB.NewMachine(machine)
}

// Modify machine info
func (a *Api) addInfo(info broadcast.ModifyMachine, r *http.Request, l *slog.Logger) (*femto.EmptyBodyResponse, error) {
	l.Info("Tried to modify machine", "host", info.Hostname, "changes", info)
	response := femto.Ok[types.Unit](types.Unit{})
	return &response, a.DB.UpdateMachine(info)
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
