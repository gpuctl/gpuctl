package webapi

import (
	"encoding/json"
	"testing"

	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	var data workstations

	jsonRep := `[{"name":"Shared","workStations":[{"name":"Workstation 1","gpus":[{"gpu_name":"NVIDIA GeForce GT 1030","gpu_brand":"GeForce","driver_ver":"535.146.02","memory_total":2048,"memory_util":0,"gpu_util":0,"memory_used":82,"fan_speed":35,"gpu_temp":31}]},{"name":"Workstation 2","gpus":[{"gpu_name":"NVIDIA TITAN Xp","gpu_brand":"Titan","driver_ver":"535.146.02","memory_total":12288,"memory_util":0,"gpu_util":0,"memory_used":83,"fan_speed":23,"gpu_temp":32},{"gpu_name":"NVIDIA TITAN Xp","gpu_brand":"Titan","driver_ver":"535.146.02","memory_total":12288,"memory_util":0,"gpu_util":0,"memory_used":83,"fan_speed":23,"gpu_temp":32}]},{"name":"Workstation 3","gpus":[{"gpu_name":"NVIDIA GeForce GT 730","gpu_brand":"GeForce","driver_ver":"470.223.02","memory_total":2001,"memory_util":0,"gpu_util":0,"memory_used":220,"fan_speed":30,"gpu_temp":27}]},{"name":"Workstation 5","gpus":[{"gpu_name":"NVIDIA TITAN Xp","gpu_brand":"Titan","driver_ver":"535.146.02","memory_total":12288,"memory_util":0,"gpu_util":0,"memory_used":83,"fan_speed":23,"gpu_temp":32},{"gpu_name":"NVIDIA TITAN Xp","gpu_brand":"Titan","driver_ver":"535.146.02","memory_total":12288,"memory_util":0,"gpu_util":0,"memory_used":83,"fan_speed":23,"gpu_temp":32}]},{"name":"Workstation 4","gpus":[{"gpu_name":"NVIDIA GeForce GT 1030","gpu_brand":"GeForce","driver_ver":"535.146.02","memory_total":2048,"memory_util":0,"gpu_util":0,"memory_used":82,"fan_speed":35,"gpu_temp":31}]},{"name":"Workstation 6","gpus":[{"gpu_name":"NVIDIA GeForce GT 730","gpu_brand":"GeForce","driver_ver":"470.223.02","memory_total":2001,"memory_util":0,"gpu_util":0,"memory_used":220,"fan_speed":30,"gpu_temp":27}]}]}]`

	err := json.Unmarshal([]byte(jsonRep), &data)

	assert.NoError(t, err)

	assert.Equal(
		t,

		workstations(workstations{
			workstationGroup{
				Name: "Shared", WorkStations: []workStationData{
					{
						Name: "Workstation 1", Gpus: []OldGPUStatSample{
							{Name: "NVIDIA GeForce GT 1030", Brand: "GeForce", DriverVersion: "535.146.02", MemoryTotal: 0x800, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 82, FanSpeed: 35, Temp: 31},
						},
					},
					{
						Name: "Workstation 2", Gpus: []OldGPUStatSample{
							{Name: "NVIDIA TITAN Xp", Brand: "Titan", DriverVersion: "535.146.02", MemoryTotal: 0x3000, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 83, FanSpeed: 23, Temp: 32},
							{Name: "NVIDIA TITAN Xp", Brand: "Titan", DriverVersion: "535.146.02", MemoryTotal: 0x3000, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 83, FanSpeed: 23, Temp: 32},
						},
					},
					{
						Name: "Workstation 3", Gpus: []OldGPUStatSample{
							{Name: "NVIDIA GeForce GT 730", Brand: "GeForce", DriverVersion: "470.223.02", MemoryTotal: 0x7d1, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 220, FanSpeed: 30, Temp: 27},
						},
					},
					{
						Name: "Workstation 5", Gpus: []OldGPUStatSample{
							{Name: "NVIDIA TITAN Xp", Brand: "Titan", DriverVersion: "535.146.02", MemoryTotal: 0x3000, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 83, FanSpeed: 23, Temp: 32},
							{Name: "NVIDIA TITAN Xp", Brand: "Titan", DriverVersion: "535.146.02", MemoryTotal: 0x3000, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 83, FanSpeed: 23, Temp: 32},
						},
					},
					{
						Name: "Workstation 4", Gpus: []OldGPUStatSample{
							{Name: "NVIDIA GeForce GT 1030", Brand: "GeForce", DriverVersion: "535.146.02", MemoryTotal: 0x800, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 82, FanSpeed: 35, Temp: 31},
						},
					},
					{
						Name: "Workstation 6", Gpus: []OldGPUStatSample{
							{Name: "NVIDIA GeForce GT 730", Brand: "GeForce", DriverVersion: "470.223.02", MemoryTotal: 0x7d1, MemoryUtilisation: 0, GPUUtilisation: 0, MemoryUsed: 220, FanSpeed: 30, Temp: 27},
						},
					},
				}}}),
		data,
	)
}

func TestToOldGPUStats(t *testing.T) {
	s := uplink.GPUStatSample{Uuid: "stuff"}
	a := ToOldGPUStats(s)
	assert.Equal(t, s.Uuid, a.Hostname)
	assert.Equal(t, uint64(1337), a.MemoryTotal)
}
