package database

import "github.com/gpuctl/gpuctl/internal/status"

// define set of operations on the database that any provider will implement
type Database interface {
	// add the info from a heartbeat packet to the data store
	// that machine may have not been seen before
	Heartbeat(packet status.GPUStatusPacket) error

	// get the latest metrics for all approved machines
	AllMachines() ([]status.GPUStatusPacket, error)
}
