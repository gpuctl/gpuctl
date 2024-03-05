package gpustats

import (
	"math/rand"

	"github.com/gpuctl/gpuctl/internal/uplink"

	"github.com/google/uuid"
)

type FakeGPU struct {
	// two random UUIDs generated on creation
	Uuids [2]uuid.UUID
}

func (fake FakeGPU) GetGPUInformation() ([]uplink.GPUInfo, error) {
	return []uplink.GPUInfo{{
		Uuid:          fake.Uuids[0],
		Name:          "GPU-inator",
		Brand:         "doofenshmirtz evil inc",
		DriverVersion: "3.141592",
		MemoryTotal:   1,
	},
		{
			Uuid:          fake.Uuids[1],
			Name:          "GPU-inator2",
			Brand:         "doofenshmirtz evil inc",
			DriverVersion: "3.141592",
			MemoryTotal:   2,
		},
	}, nil
}
func (fake FakeGPU) GetGPUStatus() ([]uplink.GPUStatSample, error) {

	return []uplink.GPUStatSample{{
		Uuid:              fake.Uuids[0],
		MemoryUtilisation: rand.Float64(),
		GPUUtilisation:    rand.Float64(),
		MemoryUsed:        rand.Float64() * 2000,
		Temp:              rand.Float64()*40 + 20,
		FanSpeed:          rand.Float64(),
		RunningProcesses:  uplink.Processes{{Pid: 123, Name: "python-inator", MemUsed: 0, Owner: "pl"}},
	},
		{
			Uuid:              fake.Uuids[1],
			MemoryUtilisation: rand.Float64() * 100,
			GPUUtilisation:    rand.Float64() * 100,
			MemoryUsed:        rand.Float64() * 2000,
			Temp:              rand.Float64()*40 + 20,
			FanSpeed:          rand.Float64(),
			RunningProcesses:  uplink.Processes{},
		},
	}, nil
}
