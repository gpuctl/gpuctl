package gpustats

import (
	"math/rand"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

type FakeGPU struct{}

func (FakeGPU) GetGPUInformation() ([]uplink.GPUInfo, error) {
	return []uplink.GPUInfo{{
		Uuid:          "some_id",
		Name:          "GPU-inator",
		Brand:         "doofenshmirtz evil inc",
		DriverVersion: "3.141592",
		MemoryTotal:   1,
	}}, nil
}
func (FakeGPU) GetGPUStatus() ([]uplink.GPUStatSample, error) {

	return []uplink.GPUStatSample{{
		Uuid:              "some_id",
		MemoryUtilisation: rand.Float64(),
		GPUUtilisation:    rand.Float64(),
		MemoryUsed:        rand.Float64() * 2000,
		Temp:              rand.Float64()*40 + 20,
		FanSpeed:          rand.Float64(),
		RunningProcesses:  uplink.Processes{{Pid: 123, Name: "python-inator", MemUsed: 0, Owner: "pl"}},
	}}, nil
}
