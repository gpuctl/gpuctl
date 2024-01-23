package gpustats

import "github.com/gpuctl/gpuctl/internal/uplink"

type GPUDataSource interface {
	GPUStats() (uplink.GPUStats, error)
}
