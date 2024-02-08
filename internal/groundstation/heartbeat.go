package groundstation

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/types"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) heartbeat(data uplink.HeartbeatReq, req *http.Request, log *slog.Logger) (*femto.Response[types.Unit], error) {
	log.Info("Received a heartbeat", "satellite", data.Hostname)

	err := gs.db.UpdateLastSeen(data.Hostname, time.Now().Unix())

	if err != nil {
		return nil, err
	}

	response := femto.Ok[types.Unit](types.Unit{})

	return &response, nil
}
