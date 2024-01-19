// defines a number of tests for a type implementing the Database interface

package database

import (
	"testing"
)

type unitTest func(t *testing.T, db *Database);

// public unit test runner takes a testing object and an anonymous functions
// that creates a new instance of the database implementation being tested
func UnitTests(t *testing.T, emptyInstance func() *Database) {
	tests := []unitTest{databaseStartsEmpty}

	for _, test := range tests {
		db := emptyInstance()
		test(t, db)
	}
}

func databaseStartsEmpty(t *testing.T, db *Database) {
	data, err := (*db).LatestData()
	if (err == nil) {
		t.Fatalf("Unexpected error: %v", err)
	}

	size := len(data)
	if (size != 0) {
		t.Fatalf("Database is not empty initially")
	}
}
