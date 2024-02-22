package groundstation

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/tunnel"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

type ErrorDB struct{}

func (edb *ErrorDB) UpdateLastSeen(host string, time int64) error {
	return nil
}

func (edb *ErrorDB) AppendDataPoint(sample uplink.GPUStatSample) error {
	return nil
}

func (edb *ErrorDB) UpdateGPUContext(host string, info uplink.GPUInfo) error {
	return nil
}

func (edb *ErrorDB) LatestData() (broadcast.Workstations, error) {
	return nil, nil
}

func (edb *ErrorDB) LastSeen() ([]broadcast.WorkstationSeen, error) {
	return nil, errors.New("database error")
}

func (edb *ErrorDB) NewMachine(machine broadcast.NewMachine) error {
	return nil
}

func (edb *ErrorDB) UpdateMachine(changes broadcast.ModifyMachine) error {
	return nil
}

func (edb *ErrorDB) Downsample(time int64) error {
	return nil
}

func (edb *ErrorDB) RemoveMachine(changes broadcast.RemoveMachine) error {
	return nil
}

func (edb *ErrorDB) AttachFile(attach broadcast.AttachFile) error {
	return nil
}

func (edb *ErrorDB) GetFile(hostname string, filename string) (broadcast.AttachFile, error) {
	var file broadcast.AttachFile
	return file, nil
}

func (edb *ErrorDB) Drop() error {
	return nil
}

func (edb *ErrorDB) ListFiles(hostname string) ([]string, error) {
	return nil, nil
}

func (edb *ErrorDB) RemoveFile(rem broadcast.RemoveFile) error {
	return nil
}

func TestPing(t *testing.T) {
	logger := slog.Default()

	// This is non-routable as per RFC 5737
	success := ping("192.0.2.1", logger)

	if success {
		t.Error("Expected ping to fail for a non-routable IP address, but it reported success")
	} else {
		t.Log("Ping processed as expected")
	}
}

func TestPingHappy(t *testing.T) {
	if testing.Short() || os.Getenv("CI") != "" {
		t.Skip("Skipping ping test in short mode or CI environment")
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	success := ping("8.8.8.8", logger)

	if !success {
		t.Error("Failed to ping google")
	}
}

func TestMonitorWithErrorDB(t *testing.T) {
	db := &ErrorDB{}
	logger := slog.Default()

	sshConfig := tunnel.Config{User: "testuser"}

	currentTime := time.Now()
	timespanForDeath := 24 * time.Hour

	err := monitor(db, currentTime, timespanForDeath, logger, sshConfig)
	if err == nil {
		t.Fatal("Expected monitor to return an error due to ErrorDB.LastSeen, but it did not")
	}
}

func TestMonitor(t *testing.T) {
	db := database.InMemory()
	logger := slog.Default()

	sshConfig := tunnel.Config{User: "testuser"}

	currentTime := time.Now()
	db.UpdateLastSeen("machineRecent", currentTime.Unix())                 // This machine should not trigger any action
	db.UpdateLastSeen("machineOld", currentTime.Add(-48*time.Hour).Unix()) // This machine should trigger actions

	timespanForDeath := 24 * time.Hour

	err := monitor(db, currentTime, timespanForDeath, logger, sshConfig)
	if err != nil {
		t.Fatalf("Monitor encountered an error: %v", err)
	}
}
