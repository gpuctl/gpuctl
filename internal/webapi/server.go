package webapi

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

type Server struct {
	mux *femto.Femto
	api *api
}

type api struct {
	db database.Database
}

func NewServer(db database.Database) *Server {
	mux := new(femto.Femto)
	api := &api{db}

	femto.OnGet(mux, "/api/stats/all", api.allstats)
	femto.OnGet(mux, "/api/stats/offline", api.HandleOfflineMachineRequest)

	return &Server{mux, api}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	s.mux.ServeHTTP(w, r)
}

// This function involves a lot of weird unwrapping
// TODO: See if we can get the database layer to do it for us
func (a *api) allstats(r *http.Request, l *slog.Logger) (workstations, error) {
	data, err := a.db.LatestData()

	if err != nil {
		return nil, err
	}

	var ws []workStationData
	for _, machine := range data {
		if len(machine.Stats) == 0 {
			continue
		}

		gpus := make([]OldGPUStatSample, 0)
		for i := range machine.Stats {
			gpus = append(gpus, zipStats(
				machine.Hostname,
				machine.GPUInfos[i],
				machine.Stats[i],
			))
		}

		ws = append(ws, workStationData{
			Name: machine.Hostname,
			Gpus: gpus,
		})
	}

	result := []workstationGroup{{Name: "Shared", WorkStations: ws}}
	return result, nil
}

// Bodge together stats and contextual data to make OldGpuStats
func zipStats(host string, info uplink.GPUInfo, stat uplink.GPUStatSample) OldGPUStatSample {
	return OldGPUStatSample{
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
