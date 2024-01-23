// package gsapi implements the groundstation API, that satellites talk to.

package gsapi

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gpuctl/gpuctl/cmd/groundstation/remote"
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
	femto.OnPost(mux, "/api/status/", remote.HandleStatusSubmission)

	return &Server{mux, gs}
}

type groundstation struct {
	lastSeen map[string]time.Time
	mu       sync.Mutex
}

func (gs *groundstation) heartbeat(data uplink.HeartbeatReq, req *http.Request, log *slog.Logger) error {
	// TODO: Pull out just the IP here.
	gs.mu.Lock()
	defer gs.mu.Unlock()

	from := data.Hostname
	prev := gs.lastSeen[from]

	now := time.Now()
	gs.lastSeen[from] = now

	log.Info("Received a heartbeat", "satellite", from, "prev_time", prev, "cur_time", now)

	return nil
}
