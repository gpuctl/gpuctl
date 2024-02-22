package database

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

type gpuInfo struct {
	host    string
	context uplink.GPUInfo
}

type inMemory struct {
	machines map[string]broadcast.ModifyMachine // maps from hostname to machine info
	infos    map[string]gpuInfo                 // maps from uuids to context info
	stats    map[string][]uplink.GPUStatSample  // maps from uuids to slices of stats, allowing tracking of multiple datapoints
	lastSeen map[string]int64                   // map from hostname to last seen time
	mu       sync.Mutex                         // mutex to protect concurrent access
}

func InMemory() Database {
	return &inMemory{
		machines: make(map[string]broadcast.ModifyMachine),
		infos:    make(map[string]gpuInfo),
		stats:    make(map[string][]uplink.GPUStatSample),
		lastSeen: make(map[string]int64),
	}
}

func (m *inMemory) AppendDataPoint(sample uplink.GPUStatSample) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if info, pres := m.infos[sample.Uuid]; !pres {
		return ErrGpuNotPresent
	} else {
		m.stats[sample.Uuid] = append(m.stats[sample.Uuid], sample)
		m.lastSeen[info.host] = time.Now().Unix()
	}

	return nil
}

func (m *inMemory) UpdateGPUContext(host string, packet uplink.GPUInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.infos[packet.Uuid] = gpuInfo{host: host, context: packet}

	// Initialize stats slice if it doesn't exist
	if _, exists := m.stats[packet.Uuid]; !exists {
		m.stats[packet.Uuid] = []uplink.GPUStatSample{}
	}
	m.lastSeen[host] = time.Now().Unix()

	return nil
}

func (m *inMemory) LatestData() (broadcast.Workstations, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var gpus = make(map[string][]broadcast.GPU)
	for uuid, info := range m.infos {
		stats, exists := m.stats[uuid]
		if !exists || len(stats) == 0 {
			continue
		}

		stat := stats[len(stats)-1]

		inUse, user := stat.RunningProcesses.Summarise()

		gpu := broadcast.GPU{
			Uuid:          uuid,
			Name:          info.context.Name,
			Brand:         info.context.Brand,
			DriverVersion: info.context.DriverVersion,
			MemoryTotal:   info.context.MemoryTotal,
			InUse:         inUse,
			User:          user,
		}

		for _, field := range reflect.VisibleFields(reflect.TypeOf(stat)) {
			if field.Name == "Uuid" {
				continue
			}
			target := reflect.ValueOf(&gpu).Elem().FieldByName(field.Name)
			if target.CanSet() {
				target.Set(reflect.ValueOf(stat).FieldByIndex(field.Index))
			}
		}

		gpus[info.host] = append(gpus[info.host], gpu)
	}

	var groups = make(map[string][]broadcast.Workstation)
	for machine, info := range m.machines {
		group := info.Group
		if group == nil || strings.TrimSpace(*group) == "" {
			fallback := DefaultGroup
			group = &fallback
		}

		workstation := broadcast.Workstation{
			Name:        machine,
			CPU:         info.CPU,
			Motherboard: info.Motherboard,
			Notes:       info.Notes,
			LastSeen:    time.Since(time.Unix(m.lastSeen[machine], 0)),
			Gpus:        gpus[machine],
		}

		groups[*group] = append(groups[*group], workstation)
	}

	var result broadcast.Workstations
	for group, machines := range groups {
		result = append(result, broadcast.Group{Name: group, Workstations: machines})
	}

	return result, nil
}

func (m *inMemory) UpdateLastSeen(host string, time int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastSeen[host] = time
	if _, found := m.machines[host]; !found {
		m.machines[host] = broadcast.ModifyMachine{Hostname: host}
	}

	return nil
}

func (m *inMemory) LastSeen() ([]broadcast.WorkstationSeen, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var seen []broadcast.WorkstationSeen
	for name, lastSeen := range m.lastSeen {
		seen = append(seen, broadcast.WorkstationSeen{Hostname: name, LastSeen: lastSeen})
	}

	return seen, nil
}

func (m *inMemory) Downsample(cutoffTime int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for uuid, samples := range m.stats {
		var oldSamples, newSamples []uplink.GPUStatSample
		for _, sample := range samples {
			if sample.Time < cutoffTime {
				oldSamples = append(oldSamples, sample)
			} else {
				newSamples = append(newSamples, sample)
			}
		}

		sort.Slice(oldSamples, func(i, j int) bool {
			return oldSamples[i].Time < oldSamples[j].Time
		})

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

		m.stats[uuid] = append(downsampled, newSamples...)
	}

	return nil
}

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
		// Addressing potential duplicate PID handling
		for _, proc := range sample.RunningProcesses {
			processesMap[proc.Pid] = proc // Overwrites if duplicate, ensuring only the latest is kept
		}
	}

	n := float64(len(samples))
	aggregatedProcesses := make([]uplink.GPUProcInfo, 0, len(processesMap))
	for _, proc := range processesMap {
		aggregatedProcesses = append(aggregatedProcesses, proc)
	}
	// Sorting to ensure deterministic order
	sort.Slice(aggregatedProcesses, func(i, j int) bool {
		return aggregatedProcesses[i].Pid < aggregatedProcesses[j].Pid
	})

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
		Time:              minTime,
		RunningProcesses:  aggregatedProcesses,
	}

	return averagedSample
}

func (m *inMemory) NewMachine(machine broadcast.NewMachine) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.lastSeen[machine.Hostname]; exists {
		return fmt.Errorf("machine with hostname %s already exists", machine.Hostname)
	}

	m.lastSeen[machine.Hostname] = time.Now().Unix()

	newMachine := broadcast.ModifyMachine{
		Hostname: machine.Hostname,
	}

	if machine.Group != nil {
		newMachine.Group = machine.Group
	}

	m.machines[machine.Hostname] = newMachine

	return nil
}

func (m *inMemory) RemoveMachine(machine broadcast.RemoveMachine) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var uuidToRemove string
	for uuid, info := range m.infos {
		if info.host == machine.Hostname {
			uuidToRemove = uuid
			break
		}
	}

	if uuidToRemove == "" {
		return fmt.Errorf("machine with hostname %s does not exist", machine.Hostname)
	}

	delete(m.lastSeen, machine.Hostname)
	delete(m.infos, uuidToRemove)
	delete(m.stats, uuidToRemove)

	return nil
}

func (m *inMemory) UpdateMachine(changes broadcast.ModifyMachine) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.lastSeen[changes.Hostname]; !exists {
		return fmt.Errorf("machine with hostname %s does not exist", changes.Hostname)
	}

	machine, exists := m.machines[changes.Hostname]
	if !exists {
		return fmt.Errorf("machine with hostname %s does not exist for update", changes.Hostname)
	}

	if changes.CPU != nil {
		machine.CPU = changes.CPU
	}
	if changes.Motherboard != nil {
		machine.Motherboard = changes.Motherboard
	}
	if changes.Notes != nil {
		machine.Notes = changes.Notes
	}
	if changes.Group != nil {
		machine.Group = changes.Group
	}

	m.machines[changes.Hostname] = machine

	return nil
}
