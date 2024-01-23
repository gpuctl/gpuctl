package gpustats

import "errors"

type GPUStatusPacket struct {
	// Contextual information
	Name          string `json:"gpu_name"`
	Brand         string `json:"gpu_brand"`
	DriverVersion string `json:"driver_ver"`
	MemoryTotal   uint64 `json:"memory_total"`

	// Temporal statistics
	MemoryUtilisation float64 `json:"memory_util"` // Percentage of memory used
	GPUUtilisation    float64 `json:"gpu_util"`    // Percentage of memory used
	MemoryUsed        float64 `json:"memory_used"`
	FanSpeed          float64 `json:"fan_speed"` // Percentage of fan speed
	Temp              float64 `json:"gpu_temp"`  // Celcius
}

type UncontextualGPUStatusPacket struct {
	MemoryUtilisation float64 `json:"memory_util"`
	GPUUtilisation    float64 `json:"gpu_util"`
	MemoryUsed        float64 `json:"memory_used"`
	FanSpeed          float64 `json:"fan_speed"`
	Temp              float64 `json:"gpu_temp"`
}

// Combine two GPUStatusPacket instances into one.
func (p GPUStatusPacket) Add(other GPUStatusPacket) (GPUStatusPacket, error) {
	if other.Name != p.Name || other.Brand != p.Brand || other.DriverVersion != p.DriverVersion || other.MemoryTotal != p.MemoryTotal {
		return GPUStatusPacket{}, errors.New("two packets with different contexts cannot be aggregated using Add, consider using UncontextualAdd")
	}

	return GPUStatusPacket{
		Name:              other.Name,
		Brand:             other.Brand,
		DriverVersion:     other.DriverVersion,
		MemoryTotal:       other.MemoryTotal,
		MemoryUtilisation: p.MemoryUtilisation + other.MemoryUtilisation,
		GPUUtilisation:    p.GPUUtilisation + other.GPUUtilisation,
		FanSpeed:          p.FanSpeed + other.FanSpeed,
		Temp:              p.Temp + other.Temp,
		MemoryUsed:        p.MemoryUsed + other.MemoryUsed,
	}, nil
}

func (p GPUStatusPacket) UncontextualAdd(other GPUStatusPacket) UncontextualGPUStatusPacket {
	return UncontextualGPUStatusPacket{
		MemoryUtilisation: p.MemoryUtilisation + other.MemoryUtilisation,
		GPUUtilisation:    p.GPUUtilisation + other.GPUUtilisation,
		FanSpeed:          p.FanSpeed + other.FanSpeed,
		Temp:              p.Temp + other.Temp,
		MemoryUsed:        p.MemoryUsed + other.MemoryUsed,
	}
}

func (p GPUStatusPacket) Scale(scalar float64) GPUStatusPacket {
	return GPUStatusPacket{
		Name:              p.Name,
		Brand:             p.Brand,
		DriverVersion:     p.DriverVersion,
		MemoryTotal:       p.MemoryTotal,
		MemoryUtilisation: p.MemoryUtilisation * scalar,
		GPUUtilisation:    p.GPUUtilisation * scalar,
		FanSpeed:          p.FanSpeed * scalar,
		Temp:              p.Temp * scalar,
		MemoryUsed:        p.MemoryUsed * scalar,
	}
}

// Identity returns the identity element of GPUStatusPacket for the UncontextualCombine operation.
func Default(name string, brand string, driverVersion string, memoryTotal uint64) GPUStatusPacket {
	return GPUStatusPacket{
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
