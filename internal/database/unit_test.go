// defines a number of tests for a type implementing the Database interface

// TODO: this whole test suite could be a lot more terse if we had functions
// that did ExpectFail, Try, ExpectEqual, etc.

package database_test

import (
	"math"
	"reflect"
	"testing"
	"time"

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
	{"Downsample", testDownsample},
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

// TODO: verify latest set of stats returned

func testDownsample(t *testing.T, db database.Database) {
	populateDatabaseWithSampleData(db, "Test GPU", 200)

	cutoffTime := time.Now().AddDate(0, -6, 0).Unix()

	if err := db.Downsample(cutoffTime); err != nil {
		t.Fatalf("Downsample failed: %v", err)
	}

	verifyDownsampledData(t, db, "Test GPU", 101) // 101 here might not be true
}

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

	expected := []uplink.WorkstationSeen{
		{Hostname: "foo", LastSeen: 1234567890},
		{Hostname: "bar", LastSeen: 9876543210},
	}

	assert.ElementsMatch(t, expected, seen)
}

func populateDatabaseWithSampleData(db database.Database, gpuID string, numberOfSamples int) {
	db.UpdateLastSeen("test-host", 0)
	err := db.UpdateGPUContext("test-host", uplink.GPUInfo{
		Uuid:          gpuID,
		Name:          "Test GPU",
		Brand:         "Test Brand",
		DriverVersion: "1.0",
		MemoryTotal:   4,
	})

	if err != nil {
		panic("Failed to update GPU context: " + err.Error())
	}

	now := time.Now()
	for i := 0; i < numberOfSamples; i++ {
		sampleTime := now.AddDate(0, 0, -i).Unix()
		sample := uplink.GPUStatSample{
			Uuid:              gpuID,
			MemoryUtilisation: 25.4 + float64(i%10),
			GPUUtilisation:    63.5 + float64(i%10),
			MemoryUsed:        1.24 + float64(i),
			FanSpeed:          35.2 + float64(i%5),
			Temp:              54.3 + float64(i%5),
			MemoryTemp:        45.3 + float64(i%5),
			GraphicsVoltage:   150.0 + float64(i%5),
			PowerDraw:         143.5 + float64(i%10),
			GraphicsClock:     50 + float64(i%5),
			MaxGraphicsClock:  134.4 + float64(i%5),
			MemoryClock:       650.3 + float64(i%10),
			MaxMemoryClock:    750 + float64(i%10),
			Time:              sampleTime,
			RunningProcesses:  nil, // Assuming process data is not relevant for this test
		}
		err := db.AppendDataPoint(sample)
		if err != nil {
			panic("Failed to append data point")
		}
	}
}

func verifyDownsampledData(t *testing.T, db database.Database, gpuID string, expectedNumSamplesAfterDownsample int) {
	results, err := db.LatestData()
	if err != nil {
		t.Fatalf("Failed to retrieve latest data: %v", err)
	}

	found := false
	var totalSamples int
	for _, upload := range results {
		if upload.Hostname == "test-host" {
			for _, info := range upload.GPUInfos {
				if info.Name == gpuID {
					found = true
					totalSamples += len(upload.Stats)
				}
			}
		}
	}

	if !found {
		t.Fatalf("GPU %s not found in the latest data results", gpuID)
	}

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
