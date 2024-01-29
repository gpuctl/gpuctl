package database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

var (
	ErrAppendNotPresent = errors.New("appending to non present machine")
)

type inMemory struct {
	stats    map[string][]uplink.GPUStatSample
	lastSeen map[string]struct{}
	mu       sync.Mutex
}

// InMemory makes a Database represented entirely in memory.
//
// This allows testing to occur without a full postgres
// server running.
func InMemory() Database {
	return &inMemory{
		stats:    make(map[string][]uplink.GPUStatSample),
		lastSeen: make(map[string]struct{}),
	}
}

// AppendDataPoint implements Database.
func (m *inMemory) AppendDataPoint(host string, packet uplink.GPUStatSample) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, pres := m.lastSeen[host]; !pres {
		return fmt.Errorf("%w: %s", ErrAppendNotPresent, host)
	}

	m.stats[host] = append(m.stats[host], packet)
	return nil
}

func (m *inMemory) UpdateGPUContext(host string, packet uplink.GPUInfo) error {
	// TODO: add functionality
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

// LatestData implements Database.
func (m *inMemory) LatestData() (map[string][]uplink.GPUStatSample, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.stats, nil
}

func (m *inMemory) UpdateLastSeen(host string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// The time can't be queried as of present, so we don't
	// care. Also, this method should take the time as an arg.
	m.lastSeen[host] = struct{}{}

	return nil
}

var _ Database = (*inMemory)(nil)
