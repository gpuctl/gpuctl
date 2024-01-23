package gpustats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	packet1 := Default("GPU1", "BrandA", "1.0", 1024)
	packet1.MemoryUtilisation = 50
	packet1.GPUUtilisation = 60

	packet2 := Default("GPU1", "BrandA", "1.0", 1024)
	packet2.MemoryUtilisation = 30
	packet2.GPUUtilisation = 40

	combinedPacket, err := packet1.Add(packet2)
	assert.NoError(t, err)
	assert.Equal(t, float64(80), combinedPacket.MemoryUtilisation)
	assert.Equal(t, float64(100), combinedPacket.GPUUtilisation)
}

func TestAddFail(t *testing.T) {
	packet1 := Default("GPU1", "BrandA", "1.0", 1024)
	packet1.MemoryUtilisation = 50
	packet1.GPUUtilisation = 60

	packet3 := Default("GPU2", "BrandB", "2.0", 2048)
	_, err := packet1.Add(packet3)
	assert.Error(t, err)
}

func TestUncontextualAdd(t *testing.T) {
	packet1 := Default("GPU1", "BrandA", "1.0", 1024)
	packet1.MemoryUtilisation = 50
	packet1.GPUUtilisation = 60

	packet2 := Default("GPU2", "BrandB", "2.0", 2048)
	packet2.MemoryUtilisation = 30
	packet2.GPUUtilisation = 40

	uncontextualSum := packet1.UncontextualAdd(packet2)
	assert.Equal(t, float64(80), uncontextualSum.MemoryUtilisation)
	assert.Equal(t, float64(100), uncontextualSum.GPUUtilisation)
}

func TestScale(t *testing.T) {
	packet := Default("GPU1", "BrandA", "1.0", 1024)
	packet.MemoryUtilisation = 50

	scaledPacket := packet.Scale(2)
	assert.Equal(t, float64(100), scaledPacket.MemoryUtilisation)
}

func TestDefault(t *testing.T) {
	packet := Default("GPU1", "BrandA", "1.0", 1024)
	assert.Equal(t, "GPU1", packet.Name)
	assert.Equal(t, float64(0), packet.MemoryUtilisation)
}
