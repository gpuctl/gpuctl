package webapi

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

type Server struct {
	mux *femto.Femto
	api *api
}

type api struct {
	// ???
}

func NewServer() *Server {
	mux := new(femto.Femto)
	api := &api{}

	femto.OnGet[workstations](mux, "/api/stats/all", api.allstats)

	return &Server{mux, api}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

var dummyData workstations = workstations{
	workstationGroup{
		Name: "Shared", WorkStations: []workStationData{
			{
				Name: "Workstation 1", Gpus: []uplink.GPUStats{
					{Name: "NVIDIA GeForce GT 1030", Brand: "GeForce", DriverVersion: "535.146.02", MemoryTotal: 0x800, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 82, FanSpeed: 35, Temp: 31},
				},
			},
			{
				Name: "Workstation 2", Gpus: []uplink.GPUStats{
					{Name: "NVIDIA TITAN Xp", Brand: "Titan", DriverVersion: "535.146.02", MemoryTotal: 0x3000, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 83, FanSpeed: 23, Temp: 32},
					{Name: "NVIDIA TITAN Xp", Brand: "Titan", DriverVersion: "535.146.02", MemoryTotal: 0x3000, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 83, FanSpeed: 23, Temp: 32},
				},
			},
			{
				Name: "Workstation 3", Gpus: []uplink.GPUStats{
					{Name: "NVIDIA GeForce GT 730", Brand: "GeForce", DriverVersion: "470.223.02", MemoryTotal: 0x7d1, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 220, FanSpeed: 30, Temp: 27},
				},
			},
			{
				Name: "Workstation 5", Gpus: []uplink.GPUStats{
					{Name: "NVIDIA TITAN Xp", Brand: "Titan", DriverVersion: "535.146.02", MemoryTotal: 0x3000, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 83, FanSpeed: 23, Temp: 32},
					{Name: "NVIDIA TITAN Xp", Brand: "Titan", DriverVersion: "535.146.02", MemoryTotal: 0x3000, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 83, FanSpeed: 23, Temp: 32},
				},
			},
			{
				Name: "Workstation 4", Gpus: []uplink.GPUStats{
					{Name: "NVIDIA GeForce GT 1030", Brand: "GeForce", DriverVersion: "535.146.02", MemoryTotal: 0x800, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 82, FanSpeed: 35, Temp: 31},
				},
			},
			{
				Name: "Workstation 6", Gpus: []uplink.GPUStats{
					{Name: "NVIDIA GeForce GT 730", Brand: "GeForce", DriverVersion: "470.223.02", MemoryTotal: 0x7d1, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 220, FanSpeed: 30, Temp: 27},
				},
			},
		}}}

func (a *api) allstats(r *http.Request, l *slog.Logger) (workstations, error) {
	return dummyData, nil
}
