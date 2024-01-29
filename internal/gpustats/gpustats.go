package gpustats

import (
	"errors"

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
	if r.Name != l.Name || r.Brand != l.Brand || r.DriverVersion != l.DriverVersion || r.MemoryTotal != l.MemoryTotal {
		// TODO: Expose this error publicly.
		return uplink.GPUStatSample{}, errors.New("two packets with different contexts cannot be aggregated using Add, consider using UncontextualAdd")
	}

	return uplink.GPUStatSample{
		Name:              r.Name,
		Brand:             r.Brand,
		DriverVersion:     r.DriverVersion,
		MemoryTotal:       r.MemoryTotal,
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
	return uplink.GPUStatSample{
		Name:              s.Name,
		Brand:             s.Brand,
		DriverVersion:     s.DriverVersion,
		MemoryTotal:       s.MemoryTotal,
		MemoryUtilisation: s.MemoryUtilisation * scalar,
		GPUUtilisation:    s.GPUUtilisation * scalar,
		FanSpeed:          s.FanSpeed * scalar,
		Temp:              s.Temp * scalar,
		MemoryUsed:        s.MemoryUsed * scalar,
	}
}

// Identity returns the identity element of GPUStats for the UncontextualCombine operation.
func Default(name string, brand string, driverVersion string, memoryTotal uint64) uplink.GPUStatSample {
	return uplink.GPUStatSample{
		Name:          name,
		Brand:         brand,
		DriverVersion: driverVersion,
		MemoryTotal:   memoryTotal,

		MemoryUtilisation: 0,
		GPUUtilisation:    0,
		MemoryUsed:        0,
		FanSpeed:          0,
		Temp:              0,
	}
}
