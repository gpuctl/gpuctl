// defines a number of tests for a type implementing the Database interface

// TODO: this whole test suite could be a lot more terse if we had functions
// that did ExpectFail, Try, ExpectEqual, etc.

package database_test

import (
	_ "embed"
	"encoding/base64"
	"log/slog"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/uplink"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/uploadtest.pdf
var uploadPdfBytes []byte
var uploadPdfEnc = base64.StdEncoding.EncodeToString(uploadPdfBytes)

//go:embed testdata/more.txt
var uploadTxtBytes []byte
var uploadTxtEnc = base64.StdEncoding.EncodeToString(uploadTxtBytes)

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
	{"MachineInfoStartsEmpty", machineInfoStartsEmpty},
	{"MachineInfoUpdatesWork", machineInfoUpdatesWork},
	{"AttachingFiles", attachAndGetFile},
	{"AttachingFilesNonexistentHost", attachFileToNonExistentHost},
	{"GetFilesNoExist", gettingNonExistentFile},
	{"ListFiles", listFiles},
	{"RemoveFile", removeFile},
	{"RemoveNonexistentFile", removeWrongFile},
	{"MachinesCanBeRemoved", removingMachine},
	{"MachinesWithSamplesCanBeRemoved", removingMachineAndSamples},
	{"InUseInformation", inUseInformation},
	{"RemovingMachineRemovesFiles", removingMachineRemoveFiles},
	{"AddMachineAddsMachines", addingMachines},
	{"DoesNotUpdateNonexistentMachines", doesNotUpdateNonexistentMachines},
}

// fake data for adding during tests
// TODO: update with processes when they're implemented
var fakeDataInfo = uplink.GPUInfo{
	Uuid:          uuid.MustParse("7d86d61f-acb4-a007-7535-203264c18e6a"),
	Name:          "GT 1030",
	Brand:         "NVidia",
	DriverVersion: "v1.4.5",
	MemoryTotal:   4,
}

// Two fake data samples for THE SAME gpu
var fakeDataSample = uplink.GPUStatSample{
	Uuid:              uuid.MustParse("7d86d61f-acb4-a007-7535-203264c18e6a"),
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
var fakeDataSample2 = uplink.GPUStatSample{
	Uuid:              uuid.MustParse("7d86d61f-acb4-a007-7535-203264c18e6a"),
	MemoryUtilisation: 2,
	GPUUtilisation:    6,
	MemoryUsed:        1.2,
	FanSpeed:          5.2,
	Temp:              4.3,
	MemoryTemp:        4.3,
	GraphicsVoltage:   15.0,
	PowerDraw:         4.5,
	GraphicsClock:     5,
	MaxGraphicsClock:  34.4,
	MemoryClock:       6.3,
	MaxMemoryClock:    75,
	RunningProcesses:  uplink.Processes{{Pid: 3456, Name: "python", Owner: "bob"}},
}

// helper functions for getting/checking machine info
func getMachine(groups broadcast.Workstations, host string) (bool, broadcast.Group, broadcast.Workstation) {
	for _, g := range groups {
		for _, m := range g.Workstations {
			if m.Name == host {
				return true, g, m
			}
		}
	}

	return false, broadcast.Group{}, broadcast.Workstation{}
}

// functions for approximately comparing floats and data structs
const margin float64 = 0.01

func floatsNear(a float64, b float64) bool {
	return math.Abs(a-b) < margin
}
func statsNear(target broadcast.GPU, stat uplink.GPUStatSample, context uplink.GPUInfo) bool {
	// compare uuids
	if target.Uuid != stat.Uuid {
		slog.Error("stat uuid didn't match", "was", target.Uuid, "wanted", stat.Uuid)
		return false
	}
	if target.Uuid != context.Uuid {
		slog.Error("context uuid didn't match", "was", target.Uuid, "wanted", context.Uuid)
		return false
	}

	// compare running processes
	inUse, user := stat.RunningProcesses.Summarise()
	if target.InUse != inUse {
		slog.Error("InUse didn't match", "was", target.InUse, "wanted", inUse)
		return false
	}
	if target.User != user {
		slog.Error("User didn't match", "was", target.User, "wanted", user)
		return false
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
			// we've already checked running processes
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
			if from.CanUint() {
				if from.Uint() != to.Uint() {
					slog.Error("Unsigned int comparision mismatch", "field name", field.Name, "expected", from.String(), "actual", to.String())
					return false
				}
			} else if from.CanFloat() {
				if !floatsNear(from.Float(), to.Float()) {
					slog.Error("Float comparision mismatch", "field name", field.Name, "expected", from.String(), "actual", to.String())
					return false
				}
			} else if from.Type().Kind() == reflect.String {
				if from.String() != to.String() {
					slog.Error("String comparision mismatch", "field name", field.Name, "expected", from.String(), "actual", to.String())
					return false
				}
			} else {
				slog.Error("Test case for this type not yet written", "field name", field.Name, "type", field.Type.String())
				return false
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
	err = db.UpdateLastSeen("badger", time.Now())
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

	err := db.UpdateLastSeen(fakeHost, time.Now())
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

	err := db.UpdateLastSeen(fakeHost, time.Now())
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
		t.Fatalf("'results' is the wrong length/has the wrong number of groups. Expected: 1, Was: %d", len(results))
	}

	found, group, machine := getMachine(results, fakeHost)
	gpus := machine.Gpus

	if !found {
		t.Fatalf("'results' didn't contain entry for '%s'", fakeHost)
	}
	if group.Name != database.DefaultGroup {
		t.Fatalf("No group was specified for '%s', it should be in the default group. Expected '%s', Was '%s'", fakeHost, database.DefaultGroup, group.Name)
	}
	if len(gpus) != 1 {
		t.Fatalf("gpus for '%s.%s' is the wrong length. Expected: 1, Was: %d", group.Name, fakeHost, len(gpus))
	}
	if !statsNear(gpus[0], fakeDataSample, fakeDataInfo) {
		t.Fatalf("Appended data doesn't match returned latest data. Expected: %v and %v, Got: %v", fakeDataInfo, fakeDataSample, gpus[0])
	}
}

// TODO: verify datastamp changed in the database
func multipleHeartbeats(t *testing.T, db database.Database) {
	err := db.UpdateLastSeen("otter", time.Now())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = db.UpdateLastSeen("otter", time.Now())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// TODO: verify latest set of stats returned

func testLastSeen1(t *testing.T, db database.Database) {
	host := "TestHost"
	lastSeenTime := time.Now()
	db.UpdateLastSeen(host, lastSeenTime)

	lastSeenData, err := db.LastSeen()
	if err != nil {
		t.Fatalf("LastSeen failed: %v", err)
	}

	found := false
	for _, data := range lastSeenData {
		if data.Hostname == host &&
			data.LastSeen.Round(time.Second).Equal(lastSeenTime.Round(time.Second)) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("For host %s wanted to find %#v, but data only had %#v", host, lastSeenTime, lastSeenData)
	}
}

func testLastSeen2(t *testing.T, db database.Database) {

	t1 := time.Date(1, 2, 3, 4, 5, 6, 0, time.UTC)
	t2 := time.Date(7, 8, 9, 10, 11, 12, 0, time.UTC)

	err := db.UpdateLastSeen("foo", t1)
	assert.NoError(t, err)

	err = db.UpdateLastSeen("bar", t2)
	assert.NoError(t, err)

	seen, err := db.LastSeen()
	assert.NoError(t, err)
	assert.Len(t, seen, 2)

	expected := []broadcast.WorkstationSeen{
		{Hostname: "foo", LastSeen: t1.Round(time.Second)},
		{Hostname: "bar", LastSeen: t2.Round(time.Second)},
	}

	assert.Equal(t, expected[0].Hostname, "foo")
	assert.Equal(t, expected[1].Hostname, "bar")
	assert.True(t, t1.Equal(expected[0].LastSeen))
	assert.True(t, t2.Equal(expected[1].LastSeen))
}

func testAppendDataPointMissingGPU(t *testing.T, db database.Database) {
	err := db.AppendDataPoint(uplink.GPUStatSample{Uuid: uuid.MustParse("00bb654e-1823-46ae-a26c-e884e2f00ff4")})
	assert.Error(t, err)
	assert.EqualError(t, err, database.ErrGpuNotPresent.Error())
}

// test getting data all the way to a GPU
func oneGpu(t *testing.T, db database.Database) {
	data, err := db.LatestData()
	assert.NoError(t, err)
	assert.Empty(t, data)

	err = db.UpdateLastSeen("foo", time.Now())
	assert.NoError(t, err)

	err = db.UpdateGPUContext("foo", uplink.GPUInfo{})
	assert.NoError(t, err)

	err = db.AppendDataPoint(uplink.GPUStatSample{})
	assert.NoError(t, err)

	data, err = db.LatestData()
	assert.NoError(t, err)
	assert.Len(t, data, 1)
}

// a machines info starts empty
func machineInfoStartsEmpty(t *testing.T, db database.Database) {
	fakeHost := "porcupine"

	err := db.UpdateLastSeen(fakeHost, time.Now())
	assert.NoError(t, err)

	data, err := db.LatestData()
	assert.NoError(t, err)
	found, group, machine := getMachine(data, fakeHost)

	if !found {
		t.Errorf("Couldn't find machine '%s'", fakeHost)
	}

	// assert that the group is the default
	assert.Equal(t, group.Name, database.DefaultGroup)

	// check all the optional characteristics start empty
	assert.Nil(t, machine.CPU)
	assert.Nil(t, machine.Motherboard)
	assert.Nil(t, machine.Notes)
	assert.Nil(t, machine.Owner)
}

// changes to a machine are present in the result
func machineInfoUpdatesWork(t *testing.T, db database.Database) {
	fakeHost := "porcupine"

	err := db.UpdateLastSeen(fakeHost, time.Now())
	assert.NoError(t, err)

	fakeGroup := "Personal"
	fakeCPU := "Intel 8080"
	fakeMotherboard := "Connect-a-tron"
	fakeNote := "Has a fan that is very loud!"
	fakeOwner := "Billie"
	fakeChange := broadcast.ModifyMachine{
		Hostname:    fakeHost,
		CPU:         &fakeCPU,
		Motherboard: &fakeMotherboard,
		Notes:       &fakeNote,
		Group:       &fakeGroup,
		Owner:       &fakeOwner,
	}

	err = db.UpdateMachine(fakeChange)
	assert.NoError(t, err)

	data, err := db.LatestData()
	found, group, machine := getMachine(data, fakeHost)

	if !found {
		t.Errorf("Couldn't find machine '%s'", fakeHost)
	}

	assert.NotNil(t, machine.CPU)
	assert.NotNil(t, machine.Motherboard)
	assert.NotNil(t, machine.Notes)
	assert.NotNil(t, machine.Owner)

	assert.Equal(t, *machine.CPU, fakeCPU)
	assert.Equal(t, *machine.Motherboard, fakeMotherboard)
	assert.Equal(t, *machine.Notes, fakeNote)
	assert.Equal(t, *machine.Owner, fakeOwner)
	assert.Equal(t, group.Name, fakeGroup)
}

// removing a machine removes it
func removingMachine(t *testing.T, db database.Database) {
	fakeHost := "chipmunk"

	err := db.UpdateLastSeen(fakeHost, time.Now())
	assert.NoError(t, err)

	// we should find the machine now
	data, err := db.LatestData()
	assert.NoError(t, err)
	found, _, _ := getMachine(data, fakeHost)
	if !found {
		t.Error("Didn't find machine when we expected to")
	}

	err = db.RemoveMachine(broadcast.RemoveMachine{Hostname: fakeHost})
	assert.NoError(t, err)

	// we shouldn't find the machine anymore
	data, err = db.LatestData()
	assert.NoError(t, err)
	found, _, _ = getMachine(data, fakeHost)
	if found {
		t.Logf("%v", data)
		t.Error("Found the machine when we didn't expect to")
	}
}

// removing a machine removes all of its samples
func removingMachineAndSamples(t *testing.T, db database.Database) {
	fakeHost := "yak"

	err := db.UpdateLastSeen(fakeHost, time.Now())
	assert.NoError(t, err)

	// add some data
	// I'm not going to bother checking it got added, that's done by other tests
	err = db.UpdateGPUContext(fakeHost, fakeDataInfo)
	assert.NoError(t, err)
	err = db.AppendDataPoint(fakeDataSample)
	assert.NoError(t, err)

	// we should find the machine now
	data, err := db.LatestData()
	assert.NoError(t, err)
	found, _, _ := getMachine(data, fakeHost)
	if !found {
		t.Error("Didn't find machine when we expected to")
	}

	err = db.AppendDataPoint(fakeDataSample2)
	assert.NoError(t, err)

	err = db.RemoveMachine(broadcast.RemoveMachine{Hostname: fakeHost})
	assert.NoError(t, err)

	// we shouldn't find the machine anymore
	data, err = db.LatestData()
	assert.NoError(t, err)
	found, _, _ = getMachine(data, fakeHost)
	if found {
		t.Logf("%v", data)
		t.Error("Found the machine when we didn't expect to")
	}
}

// db layer handles process information
func inUseInformation(t *testing.T, db database.Database) {
	fakeHost := "hamster"
	fakeUuid := uuid.MustParse("9adb69f0-1b1c-43ce-babe-99821d2cead0")

	err := db.UpdateLastSeen(fakeHost, time.Now())
	assert.NoError(t, err)

	context := uplink.GPUInfo{
		Uuid: fakeUuid,
		Name: "jeff",
	}
	noProcesses := uplink.GPUStatSample{
		Uuid:             fakeUuid,
		RunningProcesses: make([]uplink.GPUProcInfo, 0),
	}

	oneProcess := uplink.GPUStatSample{
		Uuid: fakeUuid,
		RunningProcesses: []uplink.GPUProcInfo{
			{
				Pid:     5678,
				Name:    "python",
				MemUsed: 45.2,
				Owner:   "jeff",
			},
		},
	}

	multipleProcesses := uplink.GPUStatSample{
		Uuid: fakeUuid,
		RunningProcesses: []uplink.GPUProcInfo{
			{
				Pid:     5678,
				Name:    "python",
				MemUsed: 6.2,
				Owner:   "brenda",
			},
			{
				Pid:     53935,
				Name:    "python",
				MemUsed: 103.7,
				Owner:   "james",
			},
		},
	}

	err = db.UpdateGPUContext(fakeHost, context)
	assert.NoError(t, err)

	// send with no process information, one process and multiple processes
	for i, stat := range []uplink.GPUStatSample{noProcesses, oneProcess, multipleProcesses} {
		slog.Info("Trying user process stat sample", "index", i)
		err = db.AppendDataPoint(stat)
		assert.NoError(t, err)
		data, err := db.LatestData()
		assert.NoError(t, err)
		found, _, machine := getMachine(data, fakeHost)
		assert.True(t, found)
		assert.Len(t, machine.Gpus, 1)
		assert.True(t, statsNear(machine.Gpus[0], stat, context))
	}
}

func attachAndGetFile(t *testing.T, db database.Database) {
	fakeHost := "chipmunk"
	fakeGroup := "Shared"
	err := db.NewMachine(broadcast.NewMachine{Hostname: fakeHost, Group: &fakeGroup})
	assert.NoError(t, err)

	payload := broadcast.AttachFile{
		Hostname:    fakeHost,
		Mime:        "application/pdf",
		Filename:    "test",
		EncodedFile: uploadPdfEnc,
	}

	// Put file in db
	err = db.AttachFile(payload)
	assert.NoError(t, err)

	// Now get file
	resp, err := db.GetFile(fakeHost, payload.Filename)
	assert.NoError(t, err)
	assert.Equal(t, uploadPdfEnc, resp.EncodedFile)
	assert.Equal(t, "application/pdf", resp.Mime)
	assert.Equal(t, fakeHost, resp.Hostname)
}

func gettingNonExistentFile(t *testing.T, db database.Database) {
	fakeHost1 := "chipmunk"
	fakeHost2 := "porcupine"
	fakeGroup := "Shared"
	err := db.NewMachine(broadcast.NewMachine{Hostname: fakeHost1, Group: &fakeGroup})
	assert.NoError(t, err)

	_, err = db.GetFile(fakeHost1, "does not eexist")
	assert.ErrorIs(t, err, database.ErrFileNotPresent)

	_, err = db.GetFile(fakeHost2, "still doesnt exist")
	assert.ErrorIs(t, err, database.ErrFileNotPresent)
}

func attachFileToNonExistentHost(t *testing.T, db database.Database) {
	payload := broadcast.AttachFile{
		Hostname:    "mystery",
		Mime:        "application/pdf",
		EncodedFile: uploadPdfEnc,
	}
	err := db.AttachFile(payload)
	assert.ErrorIs(t, err, database.ErrNoSuchMachine)
}

func listFiles(t *testing.T, db database.Database) {
	fakeHost := "chipmunk"
	fakeGroup := "Shared"
	err := db.NewMachine(broadcast.NewMachine{Hostname: fakeHost, Group: &fakeGroup})
	assert.NoError(t, err)
	pdf := broadcast.AttachFile{
		Hostname:    fakeHost,
		Mime:        "application/pdf",
		Filename:    "testpdf",
		EncodedFile: uploadPdfEnc,
	}

	txt := broadcast.AttachFile{
		Hostname:    fakeHost,
		Mime:        "text/plain",
		Filename:    "testtxt",
		EncodedFile: uploadTxtEnc,
	}

	err = db.AttachFile(pdf)
	assert.NoError(t, err)
	err = db.AttachFile(txt)
	assert.NoError(t, err)

	files, err := db.ListFiles(fakeHost)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(files))

	assert.ElementsMatch(t, files, []string{pdf.Filename, txt.Filename})
}

func removeFile(t *testing.T, db database.Database) {
	fakeHost := "chestnut"
	fakeGroup := "Shared"
	err := db.NewMachine(broadcast.NewMachine{Hostname: fakeHost, Group: &fakeGroup})
	assert.NoError(t, err)

	pdf := broadcast.AttachFile{
		Hostname:    fakeHost,
		Mime:        "application/pdf",
		Filename:    "filename",
		EncodedFile: uploadPdfEnc,
	}

	err = db.AttachFile(pdf)
	assert.NoError(t, err)

	err = db.RemoveFile(broadcast.RemoveFile{Hostname: fakeHost, Filename: pdf.Filename})
	assert.NoError(t, err)

	_, err = db.GetFile(fakeHost, pdf.Filename)
	assert.ErrorIs(t, err, database.ErrFileNotPresent)
}

func removeWrongFile(t *testing.T, db database.Database) {
	fakeHost := "real"
	fakeGroup := "Shared"
	err := db.NewMachine(broadcast.NewMachine{Hostname: fakeHost, Group: &fakeGroup})
	assert.NoError(t, err)

	err = db.RemoveFile(broadcast.RemoveFile{Hostname: "mystery", Filename: "doesnt exist"})
	assert.ErrorIs(t, err, database.ErrFileNotPresent)

	err = db.RemoveFile(broadcast.RemoveFile{Hostname: "real", Filename: "doesnt exist"})
	assert.ErrorIs(t, err, database.ErrFileNotPresent)
}

func additionRemoveOldFile(t *testing.T, db database.Database) {
	fakeHost := "chestnut"
	fakeGroup := "Shared"
	err := db.NewMachine(broadcast.NewMachine{Hostname: fakeHost, Group: &fakeGroup})
	assert.NoError(t, err)

	pdf1 := broadcast.AttachFile{
		Hostname:    fakeHost,
		Mime:        "application/pdf",
		Filename:    "filename",
		EncodedFile: uploadPdfEnc,
	}

	text := broadcast.AttachFile{
		Hostname:    fakeHost,
		Mime:        "plain/text",
		Filename:    "filename",
		EncodedFile: uploadTxtEnc,
	}

	err = db.AttachFile(pdf1)
	assert.NoError(t, err)
	err = db.AttachFile(text)
	assert.NoError(t, err)
	files, err := db.ListFiles(fakeHost)
	assert.Equal(t, 1, len(files))

	file, err := db.GetFile(fakeHost, files[0])
	assert.Equal(t, text, file)
}

func removingMachineRemoveFiles(t *testing.T, db database.Database) {
	fakeHost := "chestnut"
	fakeGroup := "Shared"
	err := db.NewMachine(broadcast.NewMachine{Hostname: fakeHost, Group: &fakeGroup})
	assert.NoError(t, err)

	pdf := broadcast.AttachFile{
		Hostname:    fakeHost,
		Mime:        "application/pdf",
		Filename:    "filename",
		EncodedFile: uploadPdfEnc,
	}

	err = db.AttachFile(pdf)
	assert.NoError(t, err)

	err = db.RemoveMachine(broadcast.RemoveMachine{Hostname: fakeHost})
	assert.NoError(t, err)

	_, err = db.GetFile(fakeHost, pdf.Filename)
	assert.ErrorIs(t, err, database.ErrFileNotPresent)
}

func addingMachines(t *testing.T, db database.Database) {
	fakeHost := "chestnut"
	fakeGroup := "someGroup"
	err := db.NewMachine(broadcast.NewMachine{Hostname: fakeHost, Group: &fakeGroup})
	assert.NoError(t, err)

	seen, err := db.LatestData()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(seen))
	assert.Equal(t, seen[0].Name, fakeGroup)
	assert.Equal(t, 1, len(seen[0].Workstations))
	assert.Equal(t, seen[0].Workstations[0].Name, fakeHost)
}

func doesNotUpdateNonexistentMachines(t *testing.T, db database.Database) {
	fakeHost := "chestnut"
	group := "disregarded"
	payload := broadcast.ModifyMachine{Hostname: fakeHost, Group: &group}

	err := db.UpdateMachine(payload)
	assert.NoError(t, err)

	seen, err := db.LastSeen()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(seen))
}
