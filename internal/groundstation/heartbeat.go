package groundstation

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) heartbeat(data uplink.HeartbeatReq, req *http.Request, log *slog.Logger) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	from := data.Hostname
	prev := gs.lastSeen[from]

	now := time.Now()
	gs.lastSeen[from] = now

	log.Info("Received a heartbeat", "satellite", from, "prev_time", prev, "cur_time", now)

	return nil
}
