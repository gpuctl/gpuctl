package database

import (
	"github.com/gpuctl/gpuctl/internal/uplink"
)

// define set of operations on the database that any provider will implement
type Database interface {
	// update the last seen time for a satellite to the current time
	UpdateLastSeen(host string) error

	// record a new data point for a satellite in the data store
	AppendDataPoint(uplink.GPUStats) error

	// get the latest metrics for all approved machines
	// returns map from hostname to slice of stats of gpus on that machine
	LatestData() (map[string][]uplink.GPUStats, error)
}
