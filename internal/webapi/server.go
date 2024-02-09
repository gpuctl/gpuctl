package webapi

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	mux *femto.Femto
	api *api
}

type api struct {
	db          database.Database
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

func NewServer(db database.Database, auth femto.Authenticator[APIAuthCredientals], onboardConf OnboardConf) *Server {
	mux := new(femto.Femto)
	api := &api{db, onboardConf}

	femto.OnGet(mux, "/api/stats/all", api.allstats)

	// Set up authentication endpoint
	femto.OnPost(mux, "/api/admin/auth", func(packet APIAuthCredientals, r *http.Request, l *slog.Logger) (authResponse, error) {
		return api.authenticate(auth, packet, r, l)
	})

	// Authenticated API endpoints
	femto.OnPost(mux, "/api/admin/add_workstation", femto.AuthWrapPost(auth, femto.WrapPostFunc(api.addMachine)))
	femto.OnPost(mux, "/api/admin/stats/modify", femto.AuthWrapPost(auth, femto.WrapPostFunc(api.modifyMachineInfo)))
	femto.OnPost(mux, "/api/admin/rm_workstation", femto.AuthWrapPost(auth, femto.WrapPostFunc(api.removeMachine)))
	femto.OnGet(mux, "/api/admin/confirm",
		femto.AuthWrapGet(auth, femto.WrapGetFunc(api.confirmAdmin)))

	return &Server{mux, api}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	s.mux.ServeHTTP(w, r)
}

// This function involves a lot of weird unwrapping
// TODO: See if we can get the database layer to do it for us
func (a *api) allstats(r *http.Request, l *slog.Logger) (broadcast.Workstations, error) {
	data, err := a.db.LatestData()

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
			gpus = append(gpus, zipStats(
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
	return result, nil
}

type authResponse struct {
	Token string `json:"token"`
}

func (a *api) authenticate(auth femto.Authenticator[APIAuthCredientals], packet APIAuthCredientals, r *http.Request, l *slog.Logger) (authResponse, error) {
	// Check if credientals are correct
	tok, err := auth.CreateToken(packet)
	return authResponse{tok}, err
}

// TODO
func (a *api) addMachine(add broadcast.AddMachineInfo, r *http.Request, l *slog.Logger) error {
	err := a.onboard(broadcast.OnboardReq{Hostname: add.Hostname}, r, l)
	if err != nil {
		return err
	}
	modify := broadcast.ModifyMachine{Hostname: add.Hostname, Group: &add.Group}
	err = a.modifyMachineInfo(modify, r, l)
	return err
}

func (a *api) removeMachine(rm broadcast.RemoveMachineInfo, r *http.Request, l *slog.Logger) error {
	return a.deboard(rm, r, l)
}

func (a *api) confirmAdmin(r *http.Request, l *slog.Logger) error {
	return nil
}

// TODO
func (a *api) modifyMachineInfo(info broadcast.ModifyMachine, r *http.Request, l *slog.Logger) error {
	return a.db.UpdateMachine(info)
}

// Bodge together stats and contextual data to make OldGpuStats
func zipStats(host string, info uplink.GPUInfo, stat uplink.GPUStatSample) broadcast.OldGPUStatSample {
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
