package status

type GPUStatusPacket struct {
	// Contextual information
	Name          string `json:"gpu_name"`
	Brand         string `json:"gpu_brand"`
	DriverVersion string `json:"driver_ver"`
	MemoryTotal   uint64 `json:"memory_total"`

	// Temporal statistics
	MemoryUtilisation uint64 `json:"memory_util"` // Percentage of memory used
	GPUUtilisation    uint64 `json:"gpu_util"`    // Percentage of memory used
	MemoryUsed        uint64 `json:"memory_used"`
	FanSpeed          uint64 `json:"fan_speed"` // Percentage of fan speed
	Temp              int64  `json:"gpu_temp"`  // Celcius
}
