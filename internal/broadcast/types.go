// types for webapi<->frontend communication
package broadcast

import (
	"github.com/gpuctl/gpuctl/internal/uplink"
)

// These types need to be kept in sync with `frontend/src/Data.tsx`

type NewMachine struct {
	Hostname string  `json:"hostname"`
	Group    *string `json:"group"`
}

type RemoveMachine struct {
	Hostname string `json:"hostname"`
}

type ModifyMachine struct {
	Hostname    string  `json:"hostname"`
	CPU         *string `json:"cpu"`         // nullable - means no change
	Motherboard *string `json:"motherboard"` // nullable - means no change
	Notes       *string `json:"notes"`       // nullable - means no change
	Group       *string `json:"group"`       // nullable - means no change
}

type Workstations []WorkstationGroup

type WorkstationGroup struct {
	Name         string            `json:"name"`
	WorkStations []WorkstationData `json:"workstations"`
}

// TODO: HACK: this just uses our old GPU stat packet
type WorkstationData struct {
	Name string             `json:"name"`
	Gpus []OldGPUStatSample `json:"gpus"`
}

// TODO: HACK: this is very jerry-rigged. Delete whenever we have completely changed to the new gpu data
func ToOldGPUStats(sample uplink.GPUStatSample) OldGPUStatSample {
	return OldGPUStatSample{
		Hostname:          sample.Uuid,
		Name:              "dummy name",
		Brand:             "Dummy brand",
		DriverVersion:     "dummy driver",
		MemoryTotal:       1337,
		Uuid:              sample.Uuid,
		MemoryUtilisation: sample.MemoryUtilisation,
		GPUUtilisation:    sample.GPUUtilisation,
		MemoryUsed:        sample.MemoryUsed,
		FanSpeed:          sample.FanSpeed,
		Temp:              sample.Temp,
		MemoryTemp:        sample.MemoryTemp,
		GraphicsVoltage:   sample.GraphicsVoltage,
		PowerDraw:         sample.PowerDraw,
		GraphicsClock:     sample.GraphicsClock,
		MaxGraphicsClock:  sample.MaxGraphicsClock,
		MemoryClock:       sample.MemoryClock,
		MaxMemoryClock:    sample.MaxGraphicsClock,
	}
}

type OldGPUStatSample struct {
	Hostname          string  `json:"hostname"`
	Name              string  `json:"gpu_name"`
	Brand             string  `json:"gpu_brand"`
	DriverVersion     string  `json:"driver_ver"`
	MemoryTotal       uint64  `json:"memory_total"`
	Uuid              string  `json:"uuid"`
	MemoryUtilisation float64 `json:"memory_util"`        // Percentage of memory used
	GPUUtilisation    float64 `json:"gpu_util"`           // Percentage of memory used
	MemoryUsed        float64 `json:"memory_used"`        // In megabytes
	FanSpeed          float64 `json:"fan_speed"`          // Percentage of fan speed
	Temp              float64 `json:"gpu_temp"`           // Celcius
	MemoryTemp        float64 `json:"memory_temp"`        // Celcius
	GraphicsVoltage   float64 `json:"graphics_voltage"`   // Volts
	PowerDraw         float64 `json:"power_draw"`         // Watts
	GraphicsClock     float64 `json:"graphics_clock"`     // Mhz
	MaxGraphicsClock  float64 `json:"max_graphics_clock"` // Mhz
	MemoryClock       float64 `json:"memory_clock"`       // Mhz
	MaxMemoryClock    float64 `json:"max_memory_clock"`   // Mhz
}

type OnboardReq struct {
	Hostname string `json:"hostname"`
}

type RemoveMachineInfo struct {
	Hostname string `json:"hostname"`
}

type AddMachineInfo struct {
	Hostname string `json:"hostname"`
	Group    string `json:"group"`
}

type ModifyInfo struct {
	Hostname    string  `json:"hostname"`
	Cpu         *string `json:"cpu"`
	Motherboard *string `json:"motherboard"`
	Notes       *string `json:"notes"`
	Group       *string `json:"group"`
}
