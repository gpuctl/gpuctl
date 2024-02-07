// defines a number of tests for a type implementing the Database interface

// TODO: this whole test suite could be a lot more terse if we had functions
// that did ExpectFail, Try, ExpectEqual, etc.

package database_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

type unitTest struct {
	Name string
	F    func(t *testing.T, db database.Database)
}

// a list of tests that implementations of the Database interface should pass
var UnitTests = [...]unitTest{
	{"DatabaseStartsEmpty", databaseStartsEmpty},
	{"AppendingFailsIfMachineMissing", appendingFailsIfMachineMissing},
	{"AppendingFailsIfContextMissing", appendingFailsIfContextMissing},
	{"AppendedDataPointsAreSaved", appendedDataPointsAreSaved},
	{"MultipleHeartbeats", multipleHeartbeats},
	unitTest{"TestSuccessfulDrop", dropSuccess},
}

// fake data for adding during tests
// TODO: update with processes when they're implemented
var fakeDataInfo = uplink.GPUInfo{
	Uuid:          "GPU-7d86d61f-acb4-a007-7535-203264c18e6a",
	Name:          "GT 1030",
	Brand:         "NVidia",
	DriverVersion: "v1.4.5",
	MemoryTotal:   4,
}

// Two fake data samples for THE SAME gpu
var fakeDataSample = uplink.GPUStatSample{
	Uuid:              "GPU-7d86d61f-acb4-a007-7535-203264c18e6a",
	MemoryUtilisation: 25.4,
	GPUUtilisation:    63.5,
	MemoryUsed:        1.24,
	FanSpeed:          35.2,
	Temp:              54.3,
	MemoryTemp:        45.3,
	GraphicsVoltage:   150.0,
	PowerDraw:         143.5,
	GraphicsClock:     50,
	MaxGraphicsClock:  134.4,
	MemoryClock:       650.3,
	MaxMemoryClock:    750,
	RunningProcesses:  nil,
}

// functions for approximately comparing floats and data structs
const margin float64 = 0.01

func floatsNear(a float64, b float64) bool {
	return math.Abs(a-b) < margin
}
func statsNear(a uplink.GPUStatSample, b uplink.GPUStatSample) bool {
	aType := reflect.ValueOf(a)
	bType := reflect.ValueOf(b)

	for i := 0; i < aType.NumField(); i++ {
		aVal := aType.Field(i)

		if !aVal.CanFloat() {
			continue
		}

		if !floatsNear(aVal.Float(), bType.Field(i).Float()) {
			return false
		}
	}

	return true
}

func databaseStartsEmpty(t *testing.T, db database.Database) {
	data, err := db.LatestData()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size := len(data)
	if size != 0 {
		t.Fatalf("Database is not empty initially")
	}
}

func appendingFailsIfMachineMissing(t *testing.T, db database.Database) {
	err := db.AppendDataPoint(fakeDataSample)
	if err == nil {
		t.Fatalf("Error expected but none occurred")
	}

	// even if a different machine is present
	err = db.UpdateLastSeen("badger", 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = db.AppendDataPoint(fakeDataSample)
	if err == nil {
		t.Fatalf("Error expected but none occurred")
	}
}

func appendingFailsIfContextMissing(t *testing.T, db database.Database) {
	fakeHost := "rabbit"

	err := db.UpdateLastSeen(fakeHost, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	err = db.AppendDataPoint(fakeDataSample)
	if err == nil {
		t.Fatalf("Error expected but none occurred")
	}
}

func appendedDataPointsAreSaved(t *testing.T, db database.Database) {
	fakeHost := "elk"

	err := db.UpdateLastSeen(fakeHost, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = db.UpdateGPUContext(fakeHost, fakeDataInfo)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	db.AppendDataPoint(fakeDataSample)
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

	var found = false
	var gpus []uplink.GPUStatSample
	for _, machine := range results {
		if machine.Hostname == fakeHost {
			found = true
			gpus = machine.Stats
			break
		}
	}

	if !found {
		t.Fatalf("'results' didn't contain entry for '%s'", fakeHost)
	}
	if len(gpus) != 1 {
		t.Fatalf("'results[%s]' is the wrong length. Expected: 1, Was: %d", fakeHost, len(gpus))
	}
	if !statsNear(gpus[0], fakeDataSample) {
		t.Fatalf("Appended data doesn't match returned latest data. Expected: %v, Got: %v", fakeDataSample, gpus[0])
	}
}

// TODO: verify datastamp changed in the database
func multipleHeartbeats(t *testing.T, db database.Database) {
	err := db.UpdateLastSeen("otter", 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = db.UpdateLastSeen("otter", 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func dropSuccess(t *testing.T, db database.Database) {
	t.Parallel()

	err := db.Drop()
	if err != nil {
		t.Fatalf("Error dropping database: %v", err)
	}
}

// TODO: verify latest set of stats returned
