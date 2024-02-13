package database_test

import (
	"testing"

	"github.com/gpuctl/gpuctl/internal/database"
)

func TestInMemoryUnit(t *testing.T) {
	t.Parallel()

	for _, unit := range UnitTests {
		t.Run(unit.Name, func(t *testing.T) {
			db := database.InMemory()

			unit.F(t, db)
		})
	}
}

// STOP!!!
// Don't add any more tests to this file
// TestInMemoryUnit runs all the unit tests in unit_test.go
// Add your new test cases there
