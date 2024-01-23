package remote

import (
	"testing"

	"github.com/gpuctl/gpuctl/internal/status"
	"github.com/stretchr/testify/assert"
)

/* Status Object Construction */

func TestBuildStatusObject(t *testing.T) {
	validJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": "1.0", "memory_total": 4096, "memory_util": 50, "gpu_util": 50, "memory_used": 2048, "fan_speed": 70, "gpu_temp": 60}`)
	expectedPacket := status.GPUStatusPacket{
		Name:              "Test GPU",
		Brand:             "BrandX",
		DriverVersion:     "1.0",
		MemoryTotal:       4096,
		MemoryUtilisation: 50,
		GPUUtilisation:    50,
		MemoryUsed:        2048,
		FanSpeed:          70,
		Temp:              60,
	}

	packet, err := buildStatusObject(validJSON)
	assert.NoError(t, err)
	assert.Equal(t, expectedPacket, packet)

}

func TestBuildMalformedObject(t *testing.T) {
	invalidJSON := []byte(`{"gpu_name": "Test GPU", "gpu_brand": "BrandX", "driver_ver": 1.0}`)
	_, err := buildStatusObject(invalidJSON)

	assert.Error(t, err)
}
