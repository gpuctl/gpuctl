package gpustats_test

import (
	"testing"

	"github.com/gpuctl/gpuctl/internal/gpustats"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Parallel()
	packet1 := gpustats.Default("GPU1", "BrandA", "1.0", 1024)
	packet1.MemoryUtilisation = 50
	packet1.GPUUtilisation = 60

	packet2 := gpustats.Default("GPU1", "BrandA", "1.0", 1024)
	packet2.MemoryUtilisation = 30
	packet2.GPUUtilisation = 40

	combinedPacket, err := gpustats.Add(packet1, packet2)
	assert.NoError(t, err)
	assert.Equal(t, float64(80), combinedPacket.MemoryUtilisation)
	assert.Equal(t, float64(100), combinedPacket.GPUUtilisation)
}

func TestAddFail(t *testing.T) {
	t.Parallel()
	packet1 := gpustats.Default("GPU1", "BrandA", "1.0", 1024)
	packet1.MemoryUtilisation = 50
	packet1.GPUUtilisation = 60

	packet3 := gpustats.Default("GPU2", "BrandB", "2.0", 2048)
	_, err := gpustats.Add(packet1, packet3)
	assert.Error(t, err)
}

func TestUncontextualAdd(t *testing.T) {
	t.Parallel()
	packet1 := gpustats.Default("GPU1", "BrandA", "1.0", 1024)
	packet1.MemoryUtilisation = 50
	packet1.GPUUtilisation = 60

	packet2 := gpustats.Default("GPU2", "BrandB", "2.0", 2048)
	packet2.MemoryUtilisation = 30
	packet2.GPUUtilisation = 40

	uncontextualSum := gpustats.AddUncontextual(packet1, packet2)
	assert.Equal(t, float64(80), uncontextualSum.MemoryUtilisation)
	assert.Equal(t, float64(100), uncontextualSum.GPUUtilisation)
}

func TestScale(t *testing.T) {
	t.Parallel()
	packet := gpustats.Default("GPU1", "BrandA", "1.0", 1024)
	packet.MemoryUtilisation = 50

	scaledPacket := gpustats.Scale(packet, 2)
	assert.Equal(t, float64(100), scaledPacket.MemoryUtilisation)
}

func TestDefault(t *testing.T) {
	t.Parallel()
	packet := gpustats.Default("GPU1", "BrandA", "1.0", 1024)
	assert.Equal(t, "GPU1", packet.Name)
	assert.Equal(t, float64(0), packet.MemoryUtilisation)
}
