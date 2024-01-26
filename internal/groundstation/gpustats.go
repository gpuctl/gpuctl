package groundstation

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) gpustats(data uplink.StatsPackage, req *http.Request, log *slog.Logger) error {
	log.Info("Got GPU stats", "satellite", data.Hostname, "stats", data.Stats)

	err := gs.db.UpdateLastSeen(data.Hostname)
	if err != nil {
		return err
	}
	return gs.db.AppendDataPoint(data.Hostname, data.Stats)
}
