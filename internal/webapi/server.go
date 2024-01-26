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

	return &Server{mux, api}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (a *api) allstats(r *http.Request, l *slog.Logger) (workstations, error) {
	data, err := a.db.LatestData()

	if err != nil {
		return nil, err
	}

	var workstations []workStationData

	for name, samples := range data {
		if len(samples) == 0 {
			continue
		}

		mostRecent := samples[len(samples)-1]

		workstations = append(workstations,
			workStationData{Name: name, Gpus: []uplink.GPUStats{mostRecent}},
		)
	}

	return []workstationGroup{
		{Name: "Shared", WorkStations: workstations},
	}, nil
}
