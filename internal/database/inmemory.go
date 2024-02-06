package database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

var (
	ErrMachineNotPresent = errors.New("adding gpu to non present machine")
	ErrGpuNotPresent     = errors.New("appending to non present gpu")
)

// the info map also needs to carry around the hostname so we can link it to a
// machine when we reconstruct a GPUStatsUpload
type gpuInfo struct {
	host    string
	context uplink.GPUInfo
}
type inMemory struct {
	// maps from uuids to context info and the latest stat
	infos map[string]gpuInfo
	stats map[string]uplink.GPUStatSample
	// map from hostname to last seen time
	lastSeen map[string]struct{}
	mu       sync.Mutex
}

// InMemory makes a Database represented entirely in memory.
//
// This allows testing to occur without a full postgres
// server running.
func InMemory() Database {
	return &inMemory{
		infos:    make(map[string]gpuInfo),
		stats:    make(map[string]uplink.GPUStatSample),
		lastSeen: make(map[string]struct{}),
	}
}

// AppendDataPoint implements Database.
func (m *inMemory) AppendDataPoint(sample uplink.GPUStatSample) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, pres := m.infos[sample.Uuid]; !pres {
		return fmt.Errorf("%w: %s", ErrGpuNotPresent, sample.Uuid)
	}

	m.stats[sample.Uuid] = sample
	return nil
}

func (m *inMemory) UpdateGPUContext(host string, packet uplink.GPUInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, pres := m.lastSeen[host]; !pres {
		return fmt.Errorf("%w: %s", ErrMachineNotPresent, packet.Uuid)
	}

	m.infos[packet.Uuid] = gpuInfo{host, packet}
	return nil
}

// LatestData implements Database.
func (m *inMemory) LatestData() ([]uplink.GpuStatsUpload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// build up temporary map, then flatten into array
	// we don't set the hostname field
	var grouped = make(map[string]uplink.GpuStatsUpload, 0)
	for _, stat := range m.stats {
		// get the contextual info and hostname
		info := m.infos[stat.Uuid]
		hostname := info.host

		// get the partial list of info for this machine and create a
		// new entry in the map if not already present
		_, present := grouped[hostname]
		if !present {
			grouped[hostname] = uplink.GpuStatsUpload{
				GPUInfos: make([]uplink.GPUInfo, 0),
				Stats:    make([]uplink.GPUStatSample, 0),
			}
		}

		// this will definately work now we know an entry is present
		old := grouped[hostname]
		grouped[hostname] = uplink.GpuStatsUpload{
			GPUInfos: append(old.GPUInfos, info.context),
			Stats:    append(old.Stats, stat),
		}
	}

	// flatten map to list
	var result = make([]uplink.GpuStatsUpload, 0)
	for host, structs := range grouped {
		result = append(result, uplink.GpuStatsUpload{
			Hostname: host,
			GPUInfos: structs.GPUInfos,
			Stats:    structs.Stats,
		})
	}

	return result, nil
}

func (m *inMemory) UpdateLastSeen(host string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// The time can't be queried as of present, so we don't
	// care. Also, this method should take the time as an arg.
	m.lastSeen[host] = struct{}{}

	return nil
}

func (m *inMemory) Drop() error {
	m.mu.Lock()

	m = nil
	return nil
}
