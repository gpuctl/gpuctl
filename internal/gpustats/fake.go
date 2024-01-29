package gpustats

import (
	"math/rand"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

type FakeGPU struct{}

func (FakeGPU) GPUStats() ([]uplink.GPUStatSample, error) {

	return []uplink.GPUStatSample{{
		/* TODO: delete these
		Name:          "GPU-inator",
		Brand:         "doofenshmirtz evil inc",
		DriverVersion: "3.141592",
		MemoryTotal:   1,
		*/

		MemoryUtilisation: rand.Float64(),
		GPUUtilisation:    rand.Float64(),
		MemoryUsed:        rand.Float64() * 2000,
		Temp:              rand.Float64()*40 + 20,
		FanSpeed:          rand.Float64(),
	}}, nil
}
