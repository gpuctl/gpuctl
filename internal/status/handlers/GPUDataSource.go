package handlers

import (
	"github.com/gpuctl/gpuctl/internal/status"
)

type GPUDataSource interface {
    GetGPUStatus() (status.GPUStatusPacket, error)
}
