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
			db, err := database.Postgres(url)
			if err != nil {
				t.Fatalf("Failed to open database: %v", err)
			}

			// when the test completes, drop all tables and close db
			t.Cleanup(func() {
				err = db.Drop()
				if err != nil {
					t.Logf("Got error on cleanup: %v", err)
				}
			})

			test.F(t, db)
		})
	}
}
