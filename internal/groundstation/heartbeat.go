package groundstation

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) heartbeat(data uplink.HeartbeatReq, req *http.Request, log *slog.Logger) error {
	log.Info("Received a heartbeat", "satellite", data.Hostname)

	return gs.db.UpdateLastSeen(data.Hostname)
}
