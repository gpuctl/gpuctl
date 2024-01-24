package groundstation

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) gpustats(stats uplink.GPUStats, req *http.Request, log *slog.Logger) error {
	log.Info("Got GPU stats", "stats", stats)
	return nil
}
