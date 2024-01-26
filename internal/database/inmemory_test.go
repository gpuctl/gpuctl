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

	err = db.AppendDataPoint("foo", uplink.GPUStats{
		Name: "foo-1",
	})
	assert.NoError(t, err)

	data, err = db.LatestData()
	assert.NoError(t, err)
	assert.Len(t, data, 1)
}

func TestInMemoryUnit(t *testing.T) {
	t.Parallel()

	for _, unit := range database.UnitTests {
		t.Run(unit.Name, func(t *testing.T) {
			db := database.InMemory()

			unit.F(t, db)
		})
	}
}
