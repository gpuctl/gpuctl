package database

import "github.com/gpuctl/gpuctl/internal/status"

// define set of operations on the database that any provider will implement
type Database interface {
	// update the last seen time for a satellite to the current time
	UpdateLastSeen(host String) error

	// record a new data point for a satellite in the data store
	AppendDataPoint(packet status.GPUStatusPacket) error

	// get the latest metrics for all approved machines
	LatestData() ([]status.GPUStatusPacket, error)
}
