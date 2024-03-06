package groundstation

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/types"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) heartbeat(data uplink.HeartbeatReq, req *http.Request, log *slog.Logger) (*femto.EmptyBodyResponse, error) {
	log.Info("Received a heartbeat", "satellite", data.Hostname)

	err := gs.db.UpdateLastSeen(data.Hostname, time.Now())

	if err != nil {
		return nil, err
	}

	return femto.Ok(types.Unit{})
}
