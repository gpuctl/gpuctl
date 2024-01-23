package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gpuctl/gpuctl/cmd/groundstation/config"
	"github.com/gpuctl/gpuctl/cmd/groundstation/remote"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func main() {

	configuration, err := config.GetConfiguration("config.toml")

	if err != nil {
		// TODO: Using logging library for auditing, fail soft
		fmt.Println("Error detected when determining configuration")
		return
	}

	slog.Info("Stating groundstation", "port", configuration.Server.Port)

	// TODO: Move this into fempto.
	http.HandleFunc(uplink.StatusSubmissionUrl, remote.HandleStatusSubmission)
	http.ListenAndServe(config.PortToAddress(configuration.Server.Port), nil)

	srv := NewServer()

	http.ListenAndServe(":8080", srv)

	// TODO: Make this configurable
	// err = mux.ListenAndServe(":8080")
	slog.Info("Shut down groundstation", "err", err)
}

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

	from := req.RemoteAddr
	prev := gs.lastSeen[from]

	gs.lastSeen[from] = data.Time

	log.Info("Received a heartbeat", "prev_time", prev, "cur_time", data.Time)

	return nil
}
