package database

import (
	"sync"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

type inMemory struct {
	stats map[string][]uplink.GPUStats
	mu    sync.Mutex
}

// InMemory makes a Database represented entirely in memory.
//
// This allows testing to occur without a full postgres
// server running.
func InMemory() Database {
	return &inMemory{stats: make(map[string][]uplink.GPUStats)}
}

// AppendDataPoint implements Database.
func (m *inMemory) AppendDataPoint(host string, packet uplink.GPUStats) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stats[host] = append(m.stats[host], packet)
	return nil
}

// LatestData implements Database.
func (m *inMemory) LatestData() (map[string][]uplink.GPUStats, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.stats, nil
}

func (m *inMemory) UpdateLastSeen(host string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// This can't be queried as of present, so we don't
	// care
	return nil
}

var _ Database = (*inMemory)(nil)
