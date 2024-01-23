package postgres

import (
	"testing"

	"github.com/gpuctl/gpuctl/internal/database"
)

// function that makes a new test instance of a postgres database
func emptyInstance(pool) func() *database.Database {

	return func() *database.Database {
		return nil
	}
}

// run all the database unit tests on the postgres implementation
func TestPostgres(t *testing.T) {
	database.UnitTests(t, emptyInstance())
}
