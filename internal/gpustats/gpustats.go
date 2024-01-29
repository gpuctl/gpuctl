package gpustats

import (
	"github.com/gpuctl/gpuctl/internal/uplink"
)

type UncontextualGPUStats struct {
	MemoryUtilisation float64 `json:"memory_util"`
	GPUUtilisation    float64 `json:"gpu_util"`
	MemoryUsed        float64 `json:"memory_used"`
	FanSpeed          float64 `json:"fan_speed"`
	Temp              float64 `json:"gpu_temp"`
}

// Combine two GPUStats instances into one.
func Add(l, r uplink.GPUStatSample) (uplink.GPUStatSample, error) {
	// TODO: add new additions
	return uplink.GPUStatSample{
		MemoryUtilisation: l.MemoryUtilisation + r.MemoryUtilisation,
		GPUUtilisation:    l.GPUUtilisation + r.GPUUtilisation,
		FanSpeed:          l.FanSpeed + r.FanSpeed,
		Temp:              l.Temp + r.Temp,
		MemoryUsed:        l.MemoryUsed + r.MemoryUsed,
	}, nil
}

// Combine two UncontextualGPUStats instances into one.
func AddUncontextual(l, r uplink.GPUStatSample) UncontextualGPUStats {
	return UncontextualGPUStats{
		MemoryUtilisation: l.MemoryUtilisation + r.MemoryUtilisation,
		GPUUtilisation:    l.GPUUtilisation + r.GPUUtilisation,
		FanSpeed:          l.FanSpeed + r.FanSpeed,
		Temp:              l.Temp + r.Temp,
		MemoryUsed:        l.MemoryUsed + r.MemoryUsed,
	}
}

// Scale each value in s by scalar
func Scale(s uplink.GPUStatSample, scalar float64) uplink.GPUStatSample {
	// TODO: add new data in refactor
	return uplink.GPUStatSample{
		MemoryUtilisation: s.MemoryUtilisation * scalar,
		GPUUtilisation:    s.GPUUtilisation * scalar,
		FanSpeed:          s.FanSpeed * scalar,
		Temp:              s.Temp * scalar,
		MemoryUsed:        s.MemoryUsed * scalar,
	}
}

// Identity returns the identity element of GPUStats for the UncontextualCombine operation.
func Default(name string, brand string, driverVersion string, memoryTotal uint64) uplink.GPUStatSample {
	// TODO: add new data from refactor
	return uplink.GPUStatSample{
		MemoryUtilisation: 0,
		GPUUtilisation:    0,
		MemoryUsed:        0,
		FanSpeed:          0,
		Temp:              0,
	}
}
