package uplink

// TODO: change to "/gs-api/gpu-stats"
const GPUStatsUrl = "/gs-api/status/"

type StatsPackage struct {
	Hostname string `json:"hostname"`
	Stats GPUStats `json:"stats"`
}

type GPUStats struct {
	// Contextual information
	Name          string `json:"gpu_name"`
	Brand         string `json:"gpu_brand"`
	DriverVersion string `json:"driver_ver"`
	MemoryTotal   uint64 `json:"memory_total"`

	// Temporal statistics
	MemoryUtilisation float64 `json:"memory_util"` // Percentage of memory used
	GPUUtilisation    float64 `json:"gpu_util"`    // Percentage of memory used
	MemoryUsed        float64 `json:"memory_used"` // In megabytes
	FanSpeed          float64 `json:"fan_speed"`   // Percentage of fan speed
	Temp              float64 `json:"gpu_temp"`    // Celcius
}
