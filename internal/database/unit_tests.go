// defines a number of tests for a type implementing the Database interface

package database

import (
	"testing"
)

type unitTest struct {
	Name string
	F    func(t *testing.T, db Database)
}

// a list of tests that implementations of the Database interface should pass
var UnitTests = [...]unitTest{
	{"DatabaseStartsEmpty", databaseStartsEmpty},
}

func databaseStartsEmpty(t *testing.T, db Database) {
	data, err := db.LatestData()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size := len(data)
	if size != 0 {
		t.Fatalf("Database is not empty initially")
	}
}
