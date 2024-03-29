package database

import (
	"reflect"
	"testing"
	"time"

	"github.com/gpuctl/gpuctl/internal/uplink"

	"github.com/google/uuid"
)

func TestCalculateAverage(t *testing.T) {
	gpu1 := uuid.MustParse("2cd69871-b162-4b5c-963d-0ab9cac5baaf")

	samples := []uplink.GPUStatSample{
		{
			Uuid:              gpu1,
			MemoryUtilisation: 50,
			GPUUtilisation:    25,
			MemoryUsed:        4000,
			FanSpeed:          50,
			Temp:              70,
			MemoryTemp:        65,
			GraphicsVoltage:   1.1,
			PowerDraw:         100,
			GraphicsClock:     1500,
			MaxGraphicsClock:  1800,
			MemoryClock:       700,
			MaxMemoryClock:    900,
			Time:              1625140800, // Example timestamp
			RunningProcesses: []uplink.GPUProcInfo{
				{Pid: 1, Name: "Process1", MemUsed: 250},
			},
		},
		{
			Uuid:              gpu1,
			MemoryUtilisation: 60,
			GPUUtilisation:    35,
			MemoryUsed:        5000,
			FanSpeed:          60,
			Temp:              75,
			MemoryTemp:        70,
			GraphicsVoltage:   1.2,
			PowerDraw:         120,
			GraphicsClock:     1550,
			MaxGraphicsClock:  1850,
			MemoryClock:       750,
			MaxMemoryClock:    950,
			Time:              1625227200, // A later timestamp
			RunningProcesses: []uplink.GPUProcInfo{
				{Pid: 2, Name: "Process2", MemUsed: 300},
			},
		},
	}

	expected := uplink.GPUStatSample{
		Uuid:              gpu1,
		MemoryUtilisation: 55,
		GPUUtilisation:    30,
		MemoryUsed:        4500,
		FanSpeed:          55,
		Temp:              72.5,
		MemoryTemp:        67.5,
		GraphicsVoltage:   1.15,
		PowerDraw:         110,
		GraphicsClock:     1525,
		MaxGraphicsClock:  1825,
		MemoryClock:       725,
		MaxMemoryClock:    925,
		Time:              1625140800, // The minimum of the given timestamps
		RunningProcesses: []uplink.GPUProcInfo{
			{Pid: 1, Name: "Process1", MemUsed: 250},
			{Pid: 2, Name: "Process2", MemUsed: 300},
		},
	}

	result := CalculateAverage(samples)

	// Since floating point comparison directly might lead to precision issues,
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("calculateAverage() = %v, want %v", result, expected)
	}
}

func TestAverageProcess(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		toMerge []uplink.Processes
		merged  uplink.Processes
	}{
		{"Empty", nil, nil},
		{"Single Sample",
			[]uplink.Processes{
				{{Pid: 1, Name: "foo", MemUsed: 100}, {Pid: 2, Name: "bar", MemUsed: 200}},
			},
			[]uplink.GPUProcInfo{{Pid: 1, Name: "foo", MemUsed: 100}, {Pid: 2, Name: "bar", MemUsed: 200}},
		},
		{
			"Three Separate Processes",
			[]uplink.Processes{
				{{Pid: 1, Name: "foo", MemUsed: 100}},
				{{Pid: 2, Name: "bar", MemUsed: 200}},
				{{Pid: 3, Name: "baz", MemUsed: 300}},
			},
			uplink.Processes{
				{Pid: 1, Name: "foo", MemUsed: 100},
				{Pid: 2, Name: "bar", MemUsed: 200},
				{Pid: 3, Name: "baz", MemUsed: 300},
			},
		},
		// TODO: Test same PID appearing in multiple samples (and decide it's semantics).
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			toMerge := tc.toMerge

			samples := make([]uplink.GPUStatSample, len(toMerge))
			for i, procs := range tc.toMerge {
				samples[i].RunningProcesses = procs
			}

			merged := CalculateAverage(samples).RunningProcesses

			if !reflect.DeepEqual(merged, tc.merged) {
				t.Errorf("AverageProcesses(%v) = %v, want %v", toMerge, merged, tc.merged)
			}
		})
	}

}

// specify the time we test at
// test hasn't been passing since March started?
var now = time.Date(2024, time.February, 1, 0, 0, 0, 0, time.Local)

func TestDownsample(t *testing.T) {
	db := InMemory().(*inMemory)

	cutoffTime := now.AddDate(0, -6, 0)

	gpuUUID := uuid.MustParse("96cd8554-161d-4865-9767-60c1779c57b9")
	hostName := "test-host"

	db.infos[gpuUUID] = gpuInfo{host: hostName, context: uplink.GPUInfo{Uuid: gpuUUID}}
	db.UpdateLastSeen(hostName, now)

	for i := 1; i <= 250; i++ {
		sampleTime := now.AddDate(0, 0, -i*2).Unix() // Ensuring a spread over the year
		db.stats[gpuUUID] = append(db.stats[gpuUUID], uplink.GPUStatSample{
			Uuid:              gpuUUID,
			MemoryUtilisation: float64(i % 100), // Example values, vary as needed
			GPUUtilisation:    float64(i % 100),
			MemoryUsed:        1024 + float64(i),         // Example incremental value
			FanSpeed:          50 + float64(i%50),        // Example variation
			Temp:              60 + float64(i%40),        // Example variation
			MemoryTemp:        30 + float64(i%20),        // Example variation
			GraphicsVoltage:   1.0 + float64(i%100)/1000, // Incremental variation
			PowerDraw:         200 + float64(i%50),       // Example variation
			GraphicsClock:     1000 + float64(i%500),     // Variation
			MaxGraphicsClock:  1500 + float64(i%500),     // Variation
			MemoryClock:       500 + float64(i%250),      // Variation
			MaxMemoryClock:    750 + float64(i%250),      // Variation
			Time:              sampleTime,
			RunningProcesses: []uplink.GPUProcInfo{
				{Pid: 1234, Name: "ProcessA", MemUsed: 250.0},
				{Pid: 1235, Name: "ProcessB", MemUsed: 300.0},
			},
		})
	}

	// TODO: how did we determine this value?
	// since March we've been getting 93 rather than 94
	expectedNumSamples := 94

	if err := db.Downsample(cutoffTime); err != nil {
		t.Fatalf("Downsample failed: %v", err)
	}
	if gotNumSamples := len(db.stats[gpuUUID]); gotNumSamples != expectedNumSamples {
		t.Errorf("Downsample() resulted in %d samples for %s; want %d", gotNumSamples, gpuUUID, expectedNumSamples)
	}

}

func TestDownsamplePruneMethod(t *testing.T) {

	db := InMemory().(*inMemory)

	cutoffTime := now.AddDate(0, -6, 0)
	// use all 1s uuid for test
	gpuUUID := uuid.MustParse("95a6f0b2-634c-41ab-91e2-1b9782cf8cbd")
	hostName := "test-host"

	db.infos[gpuUUID] = gpuInfo{host: hostName, context: uplink.GPUInfo{Uuid: gpuUUID}}
	db.UpdateLastSeen(hostName, now)

	for i := 1; i <= 250; i++ {
		sampleTime := now.AddDate(0, 0, -i*2).Unix() // Ensuring a spread over the year
		db.stats[gpuUUID] = append(db.stats[gpuUUID], uplink.GPUStatSample{
			Uuid:              gpuUUID,
			MemoryUtilisation: float64(i % 100), // Example values, vary as needed
			GPUUtilisation:    float64(i % 100),
			MemoryUsed:        1024 + float64(i),         // Example incremental value
			FanSpeed:          50 + float64(i%50),        // Example variation
			Temp:              60 + float64(i%40),        // Example variation
			MemoryTemp:        30 + float64(i%20),        // Example variation
			GraphicsVoltage:   1.0 + float64(i%100)/1000, // Incremental variation
			PowerDraw:         200 + float64(i%50),       // Example variation
			GraphicsClock:     1000 + float64(i%500),     // Variation
			MaxGraphicsClock:  1500 + float64(i%500),     // Variation
			MemoryClock:       500 + float64(i%250),      // Variation
			MaxMemoryClock:    750 + float64(i%250),      // Variation
			Time:              sampleTime,
			RunningProcesses: []uplink.GPUProcInfo{
				{Pid: 1234, Name: "ProcessA", MemUsed: 250.0},
				{Pid: 1235, Name: "ProcessB", MemUsed: 300.0},
			},
		})
	}

	expectedNumSamples := 94

	if err := downsampleDatabase(db, cutoffTime); err != nil {
		t.Fatalf("Downsample failed: %v", err)
	}
	if gotNumSamples := len(db.stats[gpuUUID]); gotNumSamples != expectedNumSamples {
		t.Errorf("Downsample() resulted in %d samples for %s; want %d", gotNumSamples, gpuUUID, expectedNumSamples)
	}
}
