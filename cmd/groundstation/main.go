package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/cmd/groundstation/config"
	"github.com/gpuctl/gpuctl/cmd/groundstation/remote"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func main() {
	slog.Info("Starting groundstation")

	configuration, err := config.GetConfiguration("config.toml")

	if err != nil {
		// TODO: Using logging library for auditing, fail soft
		fmt.Println("Error detected when determining configuration")
		return
	}

	// TODO: Move this into fempto.
	http.HandleFunc("/api/remote", remote.HandleStatusSubmission)
	http.ListenAndServe(config.PortToAddress(configuration.Server.Port), nil)

	mux := new(femto.Femto)
	gs := groundstation{
		lastSeen: make(map[string]time.Time),
	}
	femto.OnPost(mux, uplink.HeartbeatUrl, gs.heartbeat)

	// TODO: Make this configurable
	err = mux.ListenAndServe(":8080")
	slog.Info("Shut down groundstation", "err", err)
}

type groundstation struct {
	lastSeen map[string]time.Time
}

func (gs *groundstation) heartbeat(data uplink.HeartbeatReq, req *http.Request, log *slog.Logger) error {
	// TODO: Pull out just the IP here.
	from := req.RemoteAddr
	prev := gs.lastSeen[from]

	gs.lastSeen[from] = data.Time

	log.Info("Received a heartbeat", "prev_time", prev, "cur_time", data.Time)

	return nil
}
