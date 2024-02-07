package groundstation

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/types"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) heartbeat(data uplink.HeartbeatReq, req *http.Request, log *slog.Logger) (*femto.HTTPResponseContent[types.Unit], error) {
	log.Info("Received a heartbeat", "satellite", data.Hostname)

	// ! Yes, I am making an exception here. The response to a heartbeat is powered by femto, but should be lightweight,
	// ! and should not rely on femto's FailHandler.

	// ! Even if the heartbeat isn't recorded on our side properly, what good is it telling them that?
	return &femto.HTTPResponseContent[types.Unit]{}, gs.db.UpdateLastSeen(data.Hostname, time.Now().Unix())
}
