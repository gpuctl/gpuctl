package postgres

import (
	"os"
	"testing"

	"github.com/gpuctl/gpuctl/internal/database"
)

// run all the database unit tests on the postgres implementation
func TestPostgres(t *testing.T) {
	for _, test := range database.UnitTests {
		t.Run(test.Name, func(t *testing.T) {
			db, err := New(os.Getenv("TEST_URL"))
			if err != nil {
				t.Fatalf("Failed to open database: %v", err)
			}

			// when the test completes, drop all tables and close db
			t.Cleanup(func() {
				conn := db.(postgresConn)
				_, err = conn.db.Exec(`DROP TABLE stats;
					DROP TABLE gpus;
					DROP TABLE machines`)
				if err != nil {
					t.Logf("Got error on cleanup: %v", err)
				}
				err = conn.db.Close()
				if err != nil {
					t.Logf("Got error on cleanup: %v", err)
				}
			})

			test.F(t, db)
		})
	}
}
