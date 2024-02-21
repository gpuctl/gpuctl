// types for webapi<->frontend communication
package broadcast

import (
	"time"
)

// frontend<->web-api types
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

type AttachFile struct {
	Hostname    string `json:"hostname"`
	Mime        string `json:"mime"`
	Filename 	string `json:"filename"`
	EncodedFile string `json:"file_enc"`
}

type RemoveFile struct {
	Hostname    string `json:"hostname"`
	Filename 	string `json:"filename"`
}

// data type representing struct returned on all workstations request
type Workstations []Group

type Group struct {
	Name         string        `json:"name"` // group name
	Workstations []Workstation `json:"workstations"`
}

type Workstation struct {
	Name        string        `json:"name"`        // machine hostname
	CPU         *string       `json:"cpu"`         // cpu name (optional)
	Motherboard *string       `json:"motherboard"` // motherboard (optional)
	Notes       *string       `json:"notes"`       // general note (optional)
	LastSeen    time.Duration `json:"last_seen"`   // time since the machine was last seen
	Gpus        []GPU         `json:"gpus"`
}

type GPU struct {
	Uuid              string  `json:"uuid"`
	Name              string  `json:"gpu_name"`
	Brand             string  `json:"gpu_brand"`
	DriverVersion     string  `json:"driver_ver"`
	MemoryTotal       uint64  `json:"memory_total"`
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
	InUse             bool    `json:"in_user"`            // is this gpu being used?
	User              string  `json:"user"`               // iff it's being used, who is using this gpu
}

type OnboardReq struct {
	Hostname string `json:"hostname"`
}

type RemoveMachineInfo struct {
	Hostname string `json:"hostname"`
}

// data type returned by queries of when a workstation was last seen
type WorkstationSeen struct {
	Hostname string
	LastSeen int64
}

type DurationDeltas struct {
	Hostname string `json:"hostname"`
	Delta    int64  `json:"seconds_since"`
}
