package groundstation

import (
	"net/http"
	"sync"
	"time"

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

func NewServer() *Server {
	mux := new(femto.Femto)
	gs := &groundstation{lastSeen: make(map[string]time.Time)}

	/// Register routes.
	femto.OnPost(mux, uplink.HeartbeatUrl, gs.heartbeat)
	femto.OnPost(mux, uplink.GPUStatsUrl, gs.gpustats)

	return &Server{mux, gs}
}

type groundstation struct {
	lastSeen map[string]time.Time
	mu       sync.Mutex
}
