package groundstation

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) gpustats(data uplink.GpuStatsUpload, req *http.Request, log *slog.Logger) error {
	log.Info("Got GPU stats", "stats", data.Stats)

	// NOTE: just commented this during the big refactor -jyry
	/*
		err := gs.db.UpdateLastSeen(sample.Uuid)
		if err != nil {
			return err
		}
	*/

	for _, sample := range data.Stats {
		err := gs.db.AppendDataPoint(sample.Uuid, sample)
		if err != nil {
			return err
		}
	}
	return nil
}
