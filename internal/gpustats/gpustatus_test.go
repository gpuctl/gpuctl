package gpustats_test

import (
	"testing"

	"github.com/gpuctl/gpuctl/internal/gpustats"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	uuid1 := uuid.MustParse("b70973cb-c2c5-4fb3-b255-e537e524da5a")
	packet1 := uplink.GPUStatSample{Uuid: uuid1, MemoryUtilisation: 50, GPUUtilisation: 60}
	packet2 := uplink.GPUStatSample{Uuid: uuid1, MemoryUtilisation: 30, GPUUtilisation: 40}

	combinedPacket, err := gpustats.Add(packet1, packet2)
	assert.NoError(t, err)
	assert.Equal(t, float64(80), combinedPacket.MemoryUtilisation)
	assert.Equal(t, float64(100), combinedPacket.GPUUtilisation)
}

func TestAddFail(t *testing.T) {
	t.Parallel()

	uuid1 := uuid.MustParse("82cc3359-11f9-4041-9fbe-db784806605e")
	uuid2 := uuid.MustParse("415594ae-93c9-4cdb-8640-781be292f1d2")

	packet1 := uplink.GPUStatSample{Uuid: uuid1, MemoryUtilisation: 50, GPUUtilisation: 60}
	packet2 := uplink.GPUStatSample{Uuid: uuid2, MemoryUtilisation: 30, GPUUtilisation: 40}

	_, err := gpustats.Add(packet1, packet2)
	assert.Error(t, err)
}

func TestScale(t *testing.T) {
	t.Parallel()

	uuid1 := uuid.MustParse("03b3188f-b001-4273-8271-d437de43b729")

	packet := uplink.GPUStatSample{Uuid: uuid1, MemoryUtilisation: 50, GPUUtilisation: 60}

	scaledPacket := gpustats.Scale(packet, 2)

	assert.Equal(t, float64(100), scaledPacket.MemoryUtilisation)
	assert.Equal(t, float64(120), scaledPacket.GPUUtilisation)
}
