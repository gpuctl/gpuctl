package stats

type GPUStatsPacket struct {
	// Contextual information
	Name string
	Brand string
	DriverVersion string
	MemoryTotal uint64

	// Temporal statistics
	MemoryUsed uint64
	FanSpeed uint32
	Temp float32
}
