package database

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

// ErrMachineNotPresent error when adding GPU to a non-present machine
var ErrMachineNotPresent = errors.New("adding gpu to non present machine")

// ErrGpuNotPresent error when appending to a non-present GPU
var ErrGpuNotPresent = errors.New("appending to non present gpu")

// gpuInfo struct to carry around the hostname and GPU context
type gpuInfo struct {
	host    string
	context uplink.GPUInfo
}

// inMemory struct for the in-memory database representation
type inMemory struct {
	infos    map[string]gpuInfo                // maps from uuids to context info
	stats    map[string][]uplink.GPUStatSample // maps from uuids to slices of stats, allowing tracking of multiple datapoints
	lastSeen map[string]int64                  // map from hostname to last seen time
	mu       sync.Mutex                        // mutex for thread-safe operations
}

// InMemory creates a Database represented entirely in memory
func InMemory() Database {
	return &inMemory{
		infos:    make(map[string]gpuInfo),
		stats:    make(map[string][]uplink.GPUStatSample),
		lastSeen: make(map[string]int64),
	}
}

// AppendDataPoint appends a new GPUStatSample to the database
func (m *inMemory) AppendDataPoint(sample uplink.GPUStatSample) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if info, pres := m.infos[sample.Uuid]; !pres {
		return fmt.Errorf("%w: %s", ErrGpuNotPresent, sample.Uuid)
	} else {
		m.stats[sample.Uuid] = append(m.stats[sample.Uuid], sample)
		// Update lastSeen for the host associated with this GPU
		m.lastSeen[info.host] = time.Now().Unix()
	}

	return nil
}

// UpdateGPUContext updates or adds a GPU's context in the database
func (m *inMemory) UpdateGPUContext(host string, packet uplink.GPUInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// No need to check for host presence in lastSeen to update or add a GPU context
	m.infos[packet.Uuid] = gpuInfo{host: host, context: packet}
	// Initialize stats slice if it doesn't exist
	if _, exists := m.stats[packet.Uuid]; !exists {
		m.stats[packet.Uuid] = []uplink.GPUStatSample{}
	}
	m.lastSeen[host] = time.Now().Unix()

	return nil
}

// LatestData computes the latest data point for each GPU
func (m *inMemory) LatestData() ([]uplink.GpuStatsUpload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var uploads []uplink.GpuStatsUpload
	grouped := make(map[string]*uplink.GpuStatsUpload)

	for uuid, samples := range m.stats {
		if len(samples) == 0 {
			continue
		}
		info := m.infos[uuid]
		// Ensure we have a GpuStatsUpload for this host
		if _, exists := grouped[info.host]; !exists {
			grouped[info.host] = &uplink.GpuStatsUpload{
				Hostname: info.host,
				GPUInfos: []uplink.GPUInfo{info.context},
				Stats:    []uplink.GPUStatSample{samples[len(samples)-1]}, // Latest sample
			}
		} else {
			upload := grouped[info.host]
			upload.GPUInfos = append(upload.GPUInfos, info.context)
			upload.Stats = append(upload.Stats, samples[len(samples)-1]) // Latest sample
		}
	}

	// Convert map to slice
	for _, upload := range grouped {
		uploads = append(uploads, *upload)
	}

	return uploads, nil
}

// UpdateLastSeen updates the last seen time for a host
func (m *inMemory) UpdateLastSeen(host string, time int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastSeen[host] = time
	return nil
}

// LastSeen lists all hosts and their last seen times
func (m *inMemory) LastSeen() ([]uplink.WorkstationSeen, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var seen []uplink.WorkstationSeen
	for name, time := range m.lastSeen {
		seen = append(seen, uplink.WorkstationSeen{Hostname: name, LastSeen: time})
	}

	return seen, nil
}

func (m *inMemory) Downsample(cutoffTime int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for uuid, samples := range m.stats {
		// Separate samples into older and newer than cutoffTime
		var oldSamples, newSamples []uplink.GPUStatSample
		for _, sample := range samples {
			if sample.Time < cutoffTime {
				oldSamples = append(oldSamples, sample)
			} else {
				newSamples = append(newSamples, sample)
			}
		}

		// Sort oldSamples chronologically
		sort.Slice(oldSamples, func(i, j int) bool {
			return oldSamples[i].Time < oldSamples[j].Time
		})

		// Downsample old samples in batches of 100
		var downsampled []uplink.GPUStatSample
		for len(oldSamples) > 0 {
			batchEnd := 100
			if len(oldSamples) < 100 {
				batchEnd = len(oldSamples)
			}
			batch := oldSamples[:batchEnd]
			oldSamples = oldSamples[batchEnd:]

			averagedSample := CalculateAverage(batch)
			downsampled = append(downsampled, averagedSample)
		}

		// Combine downsampled old samples with new samples
		m.stats[uuid] = append(downsampled, newSamples...)
	}

	return nil
}

// CalculateAverage calculates the average of given GPUStatSample slices
func CalculateAverage(samples []uplink.GPUStatSample) uplink.GPUStatSample {
	if len(samples) == 0 {
		return uplink.GPUStatSample{}
	}

	var sumMemoryUtil, sumGPUUtil, sumMemoryUsed, sumFanSpeed, sumTemp, sumMemoryTemp, sumGraphicsVoltage, sumPowerDraw, sumGraphicsClock, sumMaxGraphicsClock, sumMemoryClock, sumMaxMemoryClock float64
	var minTime int64 = samples[0].Time
	processesMap := make(map[uint64]uplink.GPUProcInfo)

	for _, sample := range samples {
		sumMemoryUtil += sample.MemoryUtilisation
		sumGPUUtil += sample.GPUUtilisation
		sumMemoryUsed += sample.MemoryUsed
		sumFanSpeed += sample.FanSpeed
		sumTemp += sample.Temp
		sumMemoryTemp += sample.MemoryTemp
		sumGraphicsVoltage += sample.GraphicsVoltage
		sumPowerDraw += sample.PowerDraw
		sumGraphicsClock += sample.GraphicsClock
		sumMaxGraphicsClock += sample.MaxGraphicsClock
		sumMemoryClock += sample.MemoryClock
		sumMaxMemoryClock += sample.MaxMemoryClock

		if sample.Time < minTime {
			minTime = sample.Time
		}

		for _, proc := range sample.RunningProcesses {
			processesMap[proc.Pid] = proc
		}
	}

	n := float64(len(samples))
	aggregatedProcesses := make([]uplink.GPUProcInfo, 0, len(processesMap))
	for _, proc := range processesMap {
		aggregatedProcesses = append(aggregatedProcesses, proc)
	}

	averagedSample := uplink.GPUStatSample{
		Uuid:              samples[0].Uuid,
		MemoryUtilisation: sumMemoryUtil / n,
		GPUUtilisation:    sumGPUUtil / n,
		MemoryUsed:        sumMemoryUsed / n,
		FanSpeed:          sumFanSpeed / n,
		Temp:              sumTemp / n,
		MemoryTemp:        sumMemoryTemp / n,
		GraphicsVoltage:   sumGraphicsVoltage / n,
		PowerDraw:         sumPowerDraw / n,
		GraphicsClock:     sumGraphicsClock / n,
		MaxGraphicsClock:  sumMaxGraphicsClock / n,
		MemoryClock:       sumMemoryClock / n,
		MaxMemoryClock:    sumMaxMemoryClock / n,
		Time:              minTime, // Use the earliest time as the timestamp for the averaged sample
		RunningProcesses:  aggregatedProcesses,
	}

	return averagedSample
}

func (m *inMemory) Drop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Reset each map to a new, empty instance.
	m.infos = make(map[string]gpuInfo)
	m.stats = make(map[string][]uplink.GPUStatSample)
	m.lastSeen = make(map[string]int64)

	return nil
}
