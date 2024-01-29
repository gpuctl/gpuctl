package database

import (
	"github.com/gpuctl/gpuctl/internal/uplink"
)

// define set of operations on the database that any provider will implement
type Database interface {
	// update the last seen time for a satellite to the current time
	UpdateLastSeen(host string) error

	// record a new data point for a satellite in the data store
	// TODO: rework to take snapshot for multiple GPUs
	//       this initial version assumes all machines have a single GPU
	AppendDataPoint(host string, packet uplink.GPUStatSample) error

	// get the latest metrics for all approved machines
	// returns map from hostname to slice of stats of gpus on that machine
	LatestData() (map[string][]uplink.GPUStatSample, error)
}
