package database

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"slices"
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
	machines map[string]broadcast.ModifyMachine         // maps from hostname to machine info
	infos    map[string]gpuInfo                         // maps from uuids to context info
	stats    map[string][]uplink.GPUStatSample          // maps from uuids to slices of stats, allowing tracking of multiple datapoints
	lastSeen map[string]int64                           // map from hostname to last seen time
	files    map[string]map[string]broadcast.AttachFile // maps from hostname to attached files
	mu       sync.Mutex                                 // mutex
}

func InMemory() Database {
	return &inMemory{
		machines: make(map[string]broadcast.ModifyMachine),
		infos:    make(map[string]gpuInfo),
		stats:    make(map[string][]uplink.GPUStatSample),
		lastSeen: make(map[string]int64),
		files:    make(map[string]map[string]broadcast.AttachFile),
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

	// make mapping from machine->gpu, then make group heirarchy
	var gpus = make(map[string][]broadcast.GPU)
	for uuid, info := range m.infos {
		gpu := broadcast.GPU{
			Uuid:          uuid,
			Name:          info.context.Name,
			Brand:         info.context.Brand,
			DriverVersion: info.context.DriverVersion,
			MemoryTotal:   info.context.MemoryTotal,
		}

		// most recent stat is at the end
		stats := m.stats[uuid]
		stat := stats[len(stats)-1]

		// add on in-use info
		gpu.InUse, gpu.User = stat.RunningProcesses.Summarise()

		// reflect on stat to get all the fields
		for _, field := range reflect.VisibleFields(reflect.TypeOf(stat)) {
			// uuid already set, skip
			if field.Name == "Uuid" {
				continue
			}

			// we need a reference to gpu, not a copy so it's settable
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

		var workstation = broadcast.Workstation{
			Name:        machine,
			CPU:         info.CPU,
			Motherboard: info.Motherboard,
			Notes:       info.Notes,
			LastSeen:    time.Since(time.Unix(m.lastSeen[machine], 0)),
			Gpus:        gpus[machine],
		}

		groups[*group] = append(groups[*group], workstation)
	}

	// flatten map
	var result broadcast.Workstations = nil
	for group, machines := range groups {
		result = append(result, broadcast.Group{Name: group, Workstations: machines})
	}

	return result, nil
}

func (m *inMemory) UpdateLastSeen(host string, time int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastSeen[host] = time

	// make sure it's present in the machine list
	_, found := m.machines[host]
	if !found {
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

	// Pid -> Process
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
			// TODO: What if multiple samples have a process with the same PID?
			processesMap[proc.Pid] = proc
		}
	}

	n := float64(len(samples))

	aggregatedProcesses := make([]uplink.GPUProcInfo, 0, len(processesMap))
	for _, proc := range processesMap {
		aggregatedProcesses = append(aggregatedProcesses, proc)
	}
	// Ensure deterministic order.
	slices.SortFunc(aggregatedProcesses, func(a, b uplink.GPUProcInfo) int {
		return cmp.Compare(a.Pid, b.Pid)
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
		Time:              minTime, // Use the earliest time as the timestamp for the averaged sample
		RunningProcesses:  aggregatedProcesses,
	}

	return averagedSample
}

func (m *inMemory) NewMachine(machine broadcast.NewMachine) error {
	// TODO: add actual functionality. This was just to make the code compile
	return errors.New("NOT IMPLEMENTED FOR IN-MEMORY DB")
}

func (m *inMemory) RemoveMachine(machine broadcast.RemoveMachine) (err error) {
	// TODO: add actual functionality. This was just to make the code compile
	return errors.New("NOT IMPLEMENTED FOR IN-MEMORY DB")
}

func (m *inMemory) UpdateMachine(changes broadcast.ModifyMachine) error {
	// TODO: add actual functionality. This was just to make the code compile
	return errors.New("NOT IMPLEMENTED FOR IN-MEMORY DB")
}

func (m *inMemory) AttachFile(file broadcast.AttachFile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var zero broadcast.ModifyMachine
	if m.machines[file.Hostname] == zero {
		return fmt.Errorf("%s: %w", file.Hostname, ErrNoSuchMachine)
	}

	if m.files[file.Hostname] == nil {
		m.files[file.Hostname] = make(map[string]broadcast.AttachFile)
	}
	m.files[file.Hostname][file.Filename] = file
	return nil
}

func (m *inMemory) GetFile(hostname string, filename string) (broadcast.AttachFile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var zerofile broadcast.AttachFile
	filedir := m.files[hostname]
	if filedir == nil {
		return zerofile, fmt.Errorf("%s: %w", hostname, ErrFileNotPresent)
	}

	file := filedir[filename]
	if file == zerofile {
		return zerofile, fmt.Errorf("%s: %w", hostname, ErrFileNotPresent)
	}

	return file, nil
}

func (m *inMemory) ListFiles(hostname string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var zero broadcast.ModifyMachine
	if m.machines[hostname] == zero {
		return []string{}, ErrNoSuchMachine
	}

	if m.files[hostname] == nil {
		return []string{}, nil
	}
	res := []string{}
	for k, _ := range m.files[hostname] {
		res = append(res, k)
	}
	return res, nil
}

func (m *inMemory) RemoveFile(remove broadcast.RemoveFile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var zero broadcast.ModifyMachine
	if m.machines[remove.Hostname] == zero || m.files[remove.Hostname] == nil {
		return ErrFileNotPresent
	}
	delete(m.files[remove.Hostname], remove.Filename)
	return nil
}
