package groundstation

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/types"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func (gs *groundstation) gpustats(data uplink.GpuStatsUpload, req *http.Request, log *slog.Logger) (*femto.EmptyBodyResponse, error) {
	log.Info("Got GPU stats", "stats", data.Stats)

	err := gs.db.UpdateLastSeen(data.Hostname, time.Now())
	if err != nil {
		return nil, err
	}

	if len(data.GPUInfos) > 0 {
		err := gs.handleGPUInfo(data.Hostname, data.GPUInfos)
		if err != nil {
			return nil, err
		}
	}

	if len(data.Stats) > 0 {
		err := gs.handleGPUStatSamples(data.Hostname, data.Stats)
		if err == database.ErrGpuNotPresent {
			return &femto.EmptyBodyResponse{Status: http.StatusBadRequest}, err
		} else if err != nil {
			return nil, err
		}
	}

	return femto.Ok(types.Unit{})
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
