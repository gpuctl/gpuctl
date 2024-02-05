package groundstation

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

const (
	GPUStatsOkayResponse = "OK"
)

func (gs *groundstation) gpustats(data uplink.GpuStatsUpload, req *http.Request, log *slog.Logger) (string, error) {
	log.Info("Got GPU stats", "stats", data.Stats)

	// NOTE: just commented this during the big refactor -jyry
	err := gs.db.UpdateLastSeen(data.Hostname)
	if err != nil {
		return "", err
	}

	if len(data.GPUInfos) > 0 {
		err := gs.handleGPUInfo(data.Hostname, data.GPUInfos)
		if err != nil {
			return "", err
		}
	}

	if len(data.Stats) > 0 {
		err := gs.handleGPUStatSamples(data.Hostname, data.Stats)
		if err != nil {
			return "", err
		}
	}

	return GPUStatsOkayResponse, nil
}

func (gs *groundstation) handleGPUInfo(host string, infos []uplink.GPUInfo) error {
	for _, info := range infos {
		err := gs.db.UpdateGPUContext(host, info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *groundstation) handleGPUStatSamples(host string, stats []uplink.GPUStatSample) error {
	for _, sample := range stats {
		err := gs.db.AppendDataPoint(sample)
		if err != nil {
			return err
		}
	}
	return nil
}
