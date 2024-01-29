package gpustats

import "github.com/gpuctl/gpuctl/internal/uplink"

type GPUDataSource interface {
	GetGPUStatus() ([]uplink.GPUStatSample, error)
	GetGPUInformation() ([]uplink.GPUInfo, error)
}
