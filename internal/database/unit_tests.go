// defines a number of tests for a type implementing the Database interface

// TODO: this whole test suite could be a lot more terse if we had functions
// that did ExpectFail, Try, ExpectEqual, etc.

package database

import (
	"math"
	"testing"

	"github.com/gpuctl/gpuctl/internal/uplink"
)

type unitTest struct {
	Name string
	F    func(t *testing.T, db Database)
}

// a list of tests that implementations of the Database interface should pass
var UnitTests = [...]unitTest{
	{"DatabaseStartsEmpty", databaseStartsEmpty},
	{"AppendingFailsIfMachineMissing", appendingFailsIfMachineMissing},
	{"AppendedDataPointsAreSaved", appendedDataPointsAreSaved},
	{"MultipleHeartbeats", multipleHeartbeats},
}

// fake data for adding during tests
var fakeData = uplink.GPUStatSample{Name: "GT 1030", Brand: "NVidia",
	DriverVersion: "v1.4.5", MemoryTotal: 4, MemoryUtilisation: 25.4,
	GPUUtilisation: 63.5, MemoryUsed: 1.24, FanSpeed: 35.2, Temp: 54.3}

// functions for approximately comparing floats and data structs
const margin float64 = 0.01

func floatsNear(a float64, b float64) bool {
	return math.Abs(a-b) < margin
}
func statsNear(a uplink.GPUStatSample, b uplink.GPUStatSample) bool {
	return (a.Name == b.Name) &&
		(a.Brand == b.Brand) &&
		(a.DriverVersion == b.DriverVersion) &&
		(a.MemoryTotal == b.MemoryTotal) &&
		floatsNear(a.MemoryUtilisation, b.MemoryUtilisation) &&
		floatsNear(a.GPUUtilisation, b.GPUUtilisation) &&
		floatsNear(a.MemoryUsed, b.MemoryUsed) &&
		floatsNear(a.FanSpeed, b.FanSpeed) &&
		floatsNear(a.Temp, b.Temp)
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

func appendingFailsIfMachineMissing(t *testing.T, db Database) {
	err := db.AppendDataPoint("beaver", fakeData)
	if err == nil {
		t.Fatalf("Error expected but none occurred")
	}

	// even if a different machine is present
	err = db.UpdateLastSeen("badger")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = db.AppendDataPoint("beaver", fakeData)
	if err == nil {
		t.Fatalf("Error expected but none occurred")
	}
}

func appendedDataPointsAreSaved(t *testing.T, db Database) {
	fakeHost := "elk"

	err := db.UpdateLastSeen(fakeHost)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	db.AppendDataPoint(fakeHost, fakeData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	results, err := db.LatestData()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// check length of results and whether elk is present
	if len(results) != 1 {
		t.Fatalf("'results' is the wrong length. Expected: 1, Was: %d", len(results))
	}
	gpus, ok := results[fakeHost]
	if !ok {
		t.Fatalf("'results' didn't contain entry for '%s'", fakeHost)
	}
	if len(gpus) != 1 {
		t.Fatalf("'results[%s]' is the wrong length. Expected: 1, Was: %d", fakeHost, len(gpus))
	}
	if !statsNear(gpus[0], fakeData) {
		t.Fatalf("Appended data doesn't match returned latest data. Expected: %v, Got: %v", fakeData, gpus[0])
	}
}

// TODO: verify datastamp changed in the database
func multipleHeartbeats(t *testing.T, db Database) {
	err := db.UpdateLastSeen("otter")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = db.UpdateLastSeen("otter")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// TODO: verify latest set of stats returned
