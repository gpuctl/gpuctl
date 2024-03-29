package database_test

import (
	"os"
	"testing"

	"github.com/gpuctl/gpuctl/internal/database"
)

// run all the database unit tests on the postgres implementation
func TestPostgres(t *testing.T) {
	if testing.Short() {
		t.Skip("not connecting to postgres in short tests")
	}

	// set default value that matches github workflow
	url := os.Getenv("TEST_URL")
	if url == "" {
		url = "postgres://postgres@localhost/postgres"
	}

	for _, test := range UnitTests {
		t.Run(test.Name, func(t *testing.T) {
			dbi, err := database.Postgres(url)
			if err != nil {
				t.Fatalf("Failed to open database: %v", err)
			}

			// We want cast from the interface to the actual type, so we can call
			// the Drop method, which exists only for this tests, but not the application.
			db, ok := dbi.(database.PostgresConn)
			if !ok {
				t.Fatal("database.Postgres didn't return database.PostgresConn")
			}

			t.Cleanup(func() {
				if err := db.Drop(); err != nil {
					t.Fatal("Failed to drop database", err)
				}
			})

			test.F(t, db)
		})
	}
}

// STOP!!!
// Don't add any more tests to this file
// TestInMemoryUnit runs all the unit tests in unit_test.go
// Add your new test cases there
