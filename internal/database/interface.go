package database

import (
	"errors"
	"time"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

// Constant errors for failures in DB
var (
	ErrMachineFoundTwice = errors.New("machine found twice")
	ErrMachineNotPresent = errors.New("adding gpu to non present machine")
	ErrGpuNotPresent     = errors.New("appending to non present gpu")
	ErrNoSuchMachine     = errors.New("could not find given machine")
	ErrFileNotPresent    = errors.New("no file found")
	ErrNotImplemented    = errors.New("method not implemented")
)

// default group to give to machines with a null or empty group
const DefaultGroup string = "Shared"

// define set of operations on the database that any provider will implement
type Database interface {
	// update the last seen time for a satellite to the current time
	UpdateLastSeen(host string, time time.Time) error

	// record a new data point for a satellite in the data store
	// will error if this gpu hasn't sent a context packet yet
	AppendDataPoint(sample uplink.GPUStatSample) error

	// Update the information for the GPU contained in uplink.GPUInfo
	UpdateGPUContext(host string, info uplink.GPUInfo) error

	// get the latest metrics for all approved machines
	LatestData() (broadcast.Workstations, error)

	// get last seen online metric for all machines
	LastSeen() ([]broadcast.WorkstationSeen, error)

	// create and modify machines in the database
	NewMachine(machine broadcast.NewMachine) error
	RemoveMachine(machine broadcast.RemoveMachine) error
	UpdateMachine(changes broadcast.ModifyMachine) error

	// Downsample past certain threshold
	Downsample(time.Duration) error

	// Delete old stats past certain threshold
	DeleteOldStats(time.Duration) error

	// methods for interacting with files
	AttachFile(broadcast.AttachFile) error
	GetFile(hostname string, filename string) (broadcast.AttachFile, error)
	RemoveFile(broadcast.RemoveFile) error
	ListFiles(hostname string) ([]string, error)

	// Historical and aggregate data for graphs
	HistoricalData(hostname string) (broadcast.HistoricalData, error)
	AggregateData(days int) (broadcast.AggregateData, error)
}
