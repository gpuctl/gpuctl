package groundstation

import (
	"net/http"

	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

type Server struct {
	mux *femto.Femto
	gs  *groundstation
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func NewServer(db database.Database) *Server {
	mux := new(femto.Femto)
	gs := &groundstation{db}

	/// Register routes.
	femto.OnPost(mux, uplink.HeartbeatUrl, femto.WrapPostFunc(gs.heartbeat))
	femto.OnPost(mux, uplink.GPUStatsUrl, femto.WrapPostFunc(gs.gpustats))

	return &Server{mux, gs}
}

type groundstation struct {
	db database.Database
}
