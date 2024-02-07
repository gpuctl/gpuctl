package database

import (
	"github.com/gpuctl/gpuctl/internal/uplink"
)

// define set of operations on the database that any provider will implement
type Database interface {
	// update the last seen time for a satellite to the current time
	UpdateLastSeen(host string, time int64) error

	// record a new data point for a satellite in the data store
	// will error if this gpu hasn't sent a context packet yet
	AppendDataPoint(sample uplink.GPUStatSample) error

	// Update the information for the GPU contained in uplink.GPUInfo
	UpdateGPUContext(host string, info uplink.GPUInfo) error

	// get the latest metrics for all approved machines
	LatestData() ([]uplink.GpuStatsUpload, error)

	// get last seen online metric for all machines
	LastSeen() ([]uplink.WorkstationSeen, error)

	// Drop all tables and data in the db and close the connection
	Drop() error
}
