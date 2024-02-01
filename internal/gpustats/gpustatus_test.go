package gpustats_test

import (
	"testing"

	"github.com/gpuctl/gpuctl/internal/gpustats"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	packet1 := uplink.GPUStatSample{Uuid: "id", MemoryUtilisation: 50, GPUUtilisation: 60}
	packet2 := uplink.GPUStatSample{Uuid: "id", MemoryUtilisation: 30, GPUUtilisation: 40}

	combinedPacket, err := gpustats.Add(packet1, packet2)
	assert.NoError(t, err)
	assert.Equal(t, float64(80), combinedPacket.MemoryUtilisation)
	assert.Equal(t, float64(100), combinedPacket.GPUUtilisation)
}

func TestAddFail(t *testing.T) {
	t.Parallel()

	packet1 := uplink.GPUStatSample{Uuid: "id", MemoryUtilisation: 50, GPUUtilisation: 60}
	packet2 := uplink.GPUStatSample{Uuid: "not_id", MemoryUtilisation: 30, GPUUtilisation: 40}

	_, err := gpustats.Add(packet1, packet2)
	assert.Error(t, err)
}

func TestScale(t *testing.T) {
	t.Parallel()
	packet := uplink.GPUStatSample{Uuid: "id", MemoryUtilisation: 50, GPUUtilisation: 60}

	scaledPacket := gpustats.Scale(packet, 2)

	assert.Equal(t, float64(100), scaledPacket.MemoryUtilisation)
	assert.Equal(t, float64(120), scaledPacket.GPUUtilisation)
}
