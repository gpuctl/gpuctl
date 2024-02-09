package database_test

import (
	"testing"

	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/stretchr/testify/assert"
)

func TestInMemory(t *testing.T) {
	t.Parallel()

	db := database.InMemory()

	data, err := db.LatestData()
	assert.NoError(t, err)
	assert.Empty(t, data)

	err = db.UpdateLastSeen("foo", 0)
	assert.NoError(t, err)

	err = db.UpdateGPUContext("foo", uplink.GPUInfo{})
	assert.NoError(t, err)

	err = db.AppendDataPoint(uplink.GPUStatSample{})
	assert.NoError(t, err)

	data, err = db.LatestData()
	assert.NoError(t, err)
	assert.Len(t, data, 1)
}

func TestInMemoryUnit(t *testing.T) {
	t.Parallel()

	for _, unit := range UnitTests {
		t.Run(unit.Name, func(t *testing.T) {
			db := database.InMemory()

			unit.F(t, db)
		})
	}
}

func TestAppendDataPoint_Error(t *testing.T) {
	t.Parallel()

	db := database.InMemory()

	err := db.AppendDataPoint(uplink.GPUStatSample{Uuid: "550e8400-e29b-41d4-a716-446655440000"})
	assert.Error(t, err)
	assert.EqualError(t, err, database.ErrGpuNotPresent.Error()+": 550e8400-e29b-41d4-a716-446655440000")
}

func TestDrop(t *testing.T) {
	t.Parallel()

	db := database.InMemory()

	err := db.Drop()
	assert.NoError(t, err)
}

func TestLastSeen(t *testing.T) {
	t.Parallel()

	db := database.InMemory()

	err := db.UpdateLastSeen("foo", 1234567890)
	assert.NoError(t, err)

	err = db.UpdateLastSeen("bar", 9876543210)
	assert.NoError(t, err)

	seen, err := db.LastSeen()
	assert.NoError(t, err)
	assert.Len(t, seen, 2)

	expected := []uplink.WorkstationSeen{
		{Hostname: "foo", LastSeen: 1234567890},
		{Hostname: "bar", LastSeen: 9876543210},
	}

	assert.ElementsMatch(t, expected, seen)
}
