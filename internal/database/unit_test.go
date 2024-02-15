// defines a number of tests for a type implementing the Database interface

// TODO: this whole test suite could be a lot more terse if we had functions
// that did ExpectFail, Try, ExpectEqual, etc.

package database_test

import (
	"log/slog"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/stretchr/testify/assert"
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
	{"TestAppendDataPointMissingGPU", testAppendDataPointMissingGPU},
	{"LastSeen1", testLastSeen1},
	{"LastSeen2", testLastSeen2},
	{"OneGpu", oneGpu},
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
func statsNear(target broadcast.GPU, stat uplink.GPUStatSample, context uplink.GPUInfo) bool {
	if target.Uuid != stat.Uuid {
		slog.Error("stat uuid didn't match", "was", target.Uuid, "wanted", stat.Uuid)
		return false
	}
	if target.Uuid != context.Uuid {
		slog.Error("context uuid didn't match", "was", target.Uuid, "wanted", context.Uuid)
	}

	// compare all the other fields using reflection
	for _, compare := range []interface{}{stat, context} {
		compareV := reflect.ValueOf(compare)

		for _, field := range reflect.VisibleFields(compareV.Type()) {
			// we've already compared uuids
			if field.Name == "Uuid" {
				continue
			}
			// TODO: determine where we use time field
			if field.Name == "Time" {
				continue
			}
			// TODO: Start actually using running processes
			if field.Name == "RunningProcesses" {
				continue
			}

			// get fields from structs
			from := compareV.FieldByIndex(field.Index)
			to := reflect.ValueOf(target).FieldByName(field.Name)
			if !to.IsValid() {
				slog.Error("Couldn't get field from target struct", "field name", field.Name)
				return false
			}
			if from.Type() != to.Type() {
				slog.Error("Comparision type mismatch", "field name", field.Name, "expected", from.Type().String())
				return false
			}

			// do a different comparision based on type
			if from.CanInt() {
				if from.Int() != to.Int() {
					slog.Error("Int comparision mismatch", "field name", field.Name, "expected", from.String(), "actual", to.String())
					return false
				}
			} else if from.CanFloat() {
				if !floatsNear(from.Float(), to.Float()) {
					slog.Error("Float comparision mismatch", "field name", field.Name, "expected", from.String(), "actual", to.String())
					return false
				}
			} else {
				slog.Error("Test case for this type not yet written", "field name", field.Name, "type", field.Type.String())
			}
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
	var gpus []broadcast.GPU
	var foundGroup string
	for _, group := range results {
		for _, machine := range group.Workstations {
			if machine.Name == fakeHost {
				found = true
				gpus = machine.Gpus
				foundGroup = group.Name
				break
			}
		}
	}

	if !found {
		t.Fatalf("'results' didn't contain entry for '%s'", fakeHost)
	}
	if len(gpus) != 1 {
		t.Fatalf("gpus for '%s.%s' is the wrong length. Expected: 1, Was: %d", foundGroup, fakeHost, len(gpus))
	}
	if !statsNear(gpus[0], fakeDataSample, fakeDataInfo) {
		t.Fatalf("Appended data doesn't match returned latest data. Expected: %v and %v, Got: %v", fakeDataInfo, fakeDataSample, gpus[0])
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

// TODO: verify latest set of stats returned

func testLastSeen1(t *testing.T, db database.Database) {
	host := "TestHost"
	lastSeenTime := time.Now().Unix()
	db.UpdateLastSeen(host, lastSeenTime)

	lastSeenData, err := db.LastSeen()
	if err != nil {
		t.Fatalf("LastSeen failed: %v", err)
	}

	found := false
	for _, data := range lastSeenData {
		if data.Hostname == host && data.LastSeen == lastSeenTime {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Last seen data for host %s was not updated correctly", host)
	}
}

func testLastSeen2(t *testing.T, db database.Database) {
	err := db.UpdateLastSeen("foo", 1234567890)
	assert.NoError(t, err)

	err = db.UpdateLastSeen("bar", 9876543210)
	assert.NoError(t, err)

	seen, err := db.LastSeen()
	assert.NoError(t, err)
	assert.Len(t, seen, 2)

	expected := []broadcast.WorkstationSeen{
		{Hostname: "foo", LastSeen: 1234567890},
		{Hostname: "bar", LastSeen: 9876543210},
	}

	assert.ElementsMatch(t, expected, seen)
}

func testAppendDataPointMissingGPU(t *testing.T, db database.Database) {
	err := db.AppendDataPoint(uplink.GPUStatSample{Uuid: "bogus_uuid_blah"})
	assert.Error(t, err)
	assert.EqualError(t, err, database.ErrGpuNotPresent.Error())
}

// test getting data all the way to a GPU
func oneGpu(t *testing.T, db database.Database) {
	data, err := db.LatestData()
	assert.NoError(t, err)
	assert.Empty(t, data)

	err = db.UpdateLastSeen("foo", 0)
	assert.NoError(t, err)

	err = db.UpdateGPUContext("foo", uplink.GPUInfo{})
	assert.NoError(t, err)

	err = db.AppendDataPoint(uplink.GPUStatSample{})
	assert.NoError(t, err)

	data, err = db.LatestData()
	assert.NoError(t, err)
	assert.Len(t, data, 1)
}
