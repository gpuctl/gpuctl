package stats

import (
	"time"
)

type GPUStatsPacket struct {
	// Contextual information
	Name string
	Brand string
	DriverVersion string
	MemoryTotal uint64

	// Temporal statistics
	Timestamp time.Time
	MemoryUsed uint64
	Temp float32
}

type GPUStatGetter interface {
	currentStatus() GPUStatsPacket
}
