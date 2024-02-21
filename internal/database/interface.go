package database

import (
	"errors"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

// Constant errors for failures in DB
var (
	ErrMachineNotPresent             = errors.New("adding gpu to non present machine")
	ErrGpuNotPresent                 = errors.New("appending to non present gpu")
	ErrAddingFileToNonPresentMachine = errors.New("attaching a file to a machine that does not exist")
	ErrFileNotPresent                = errors.New("no file found")
)

// default group to give to machines with a null or empty group
const DefaultGroup string = "Shared"

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
	LatestData() (broadcast.Workstations, error)

	// get last seen online metric for all machines
	LastSeen() ([]broadcast.WorkstationSeen, error)

	// create and modify machines in the database
	NewMachine(machine broadcast.NewMachine) error
	RemoveMachine(machine broadcast.RemoveMachine) error
	UpdateMachine(changes broadcast.ModifyMachine) error

	// downsample since certain unix time
	Downsample(time int64) error

	// methods for interacting with files
	AttachFile(broadcast.AttachFile) error
	GetFile(hostname string) (broadcast.AttachFile, error)
}
