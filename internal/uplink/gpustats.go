package uplink

// TODO: change to "/gs-api/gpu-stats"
const GPUStatsUrl = "/gs-api/status/"

type GpuStatsUpload struct {
	Hostname string          `json:"hostname"`
	GPUInfos []GPUInfo       `json:"information"`
	Stats    []GPUStatSample `json:"stats"`
}

// Contextual information about the GPU
type GPUInfo struct {
	Uuid          string `json:"uuid"`
	Name          string `json:"gpu_name"`
	Brand         string `json:"gpu_brand"`
	DriverVersion string `json:"driver_ver"`
	MemoryTotal   uint64 `json:"memory_total"`
}

// Temporal statistics for a GPU
type GPUStatSample struct {
	Uuid              string    `json:"uuid"`
	MemoryUtilisation float64   `json:"memory_util"`        // Percentage of memory used
	GPUUtilisation    float64   `json:"gpu_util"`           // Percentage of memory used
	MemoryUsed        float64   `json:"memory_used"`        // In megabytes
	FanSpeed          float64   `json:"fan_speed"`          // Percentage of fan speed
	Temp              float64   `json:"gpu_temp"`           // Celcius
	MemoryTemp        float64   `json:"memory_temp"`        // Celcius
	GraphicsVoltage   float64   `json:"graphics_voltage"`   // Volts
	PowerDraw         float64   `json:"power_draw"`         // Watts
	GraphicsClock     float64   `json:"graphics_clock"`     // Mhz
	MaxGraphicsClock  float64   `json:"max_graphics_clock"` // Mhz
	MemoryClock       float64   `json:"memory_clock"`       // Mhz
	MaxMemoryClock    float64   `json:"max_memory_clock"`   // Mhz
	Time              int64     `json:"time"`
	RunningProcesses  Processes `json:"processes"` // List of processes running
}

type Processes []GPUProcInfo
type GPUProcInfo struct {
	Pid     uint64  `json:"pid"`
	Name    string  `json:"name"`
	MemUsed float64 `json:"used_memory"`
	Owner   string  `json:"owner"` // NOTE: This is the user that owns the process
}

// Summarise a running processes array
// reduces to boolean specifying whether it's in use and a user
func (procs Processes) Summarise() (inUse bool, user string) {
	user = ""
	inUse = len(procs) > 0

	if inUse {
		user = procs[0].Owner
	}
	return
}
