package groundstation

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

const (
	HeartbeatOKResponse = "OK"
)

func (gs *groundstation) heartbeat(data uplink.HeartbeatReq, req *http.Request, log *slog.Logger) (string, error) {
	log.Info("Received a heartbeat", "satellite", data.Hostname)

	return HeartbeatOKResponse, gs.db.UpdateLastSeen(data.Hostname)
}
