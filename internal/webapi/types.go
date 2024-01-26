package webapi

import (
	"github.com/gpuctl/gpuctl/internal/uplink"
)

// These types need to be kept in sync with `frontend/src/Data.tsx`

type workstations []workstationGroup

type workstationGroup struct {
	Name         string            `json:"name"`
	WorkStations []workStationData `json:"workStations"`
}

type workStationData struct {
	Name string            `json:"name"`
	Gpus []uplink.GPUStats `json:"gpus"`
}
