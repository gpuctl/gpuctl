package groundstation

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/tunnel"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

type ErrorDB struct{}

var errorDbNotImplemented = errors.New("database error: using errorDB")

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
	return nil, errorDbNotImplemented
}

func (edb *ErrorDB) NewMachine(machine broadcast.NewMachine) error {
	return nil
}

func (edb *ErrorDB) UpdateMachine(changes broadcast.ModifyMachine) error {
	return nil
}

func (edb *ErrorDB) Downsample(time time.Time) error {
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

func (edb *ErrorDB) HistoricalData(hostname string) (broadcast.HistoricalData, error) {
	return nil, nil
}

func (edb *ErrorDB) AggregateData(days int) (broadcast.AggregateData, error) {
	return broadcast.AggregateData{}, nil
}

func TestPing(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	if os.Getenv("CI") != "" {
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

func sshConfig(t *testing.T) tunnel.Config {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	signer, err := ssh.NewSignerFromKey(privateKey)
	require.NoError(t, err)

	return tunnel.Config{User: "testuser", Signer: signer}

}

func TestMonitorWithErrorDB(t *testing.T) {
	db := &ErrorDB{}
	logger := slog.Default()

	currentTime := time.Now()
	cutoffTime := currentTime.Add(-24 * time.Hour)

	sshConfig := sshConfig(t)

	err := monitor(db, cutoffTime, logger, sshConfig)
	if !errors.Is(err, errorDbNotImplemented) {
		t.Fatal("Expected monitor to return an error due to ErrorDB.LastSeen, but it did not")
	}
}

func TestMonitor(t *testing.T) {
	db := database.InMemory()

	// logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger := slog.Default()

	sshConfig := sshConfig(t)

	currentTime := time.Now()
	db.UpdateLastSeen("machineRecent", currentTime.Unix())                 // This machine should not trigger any action
	db.UpdateLastSeen("machineOld", currentTime.Add(-48*time.Hour).Unix()) // This machine should trigger actions

	cutoffTime := currentTime.Add(-24 * time.Hour)
	err := monitor(db, cutoffTime, logger, sshConfig)

	assert.NoError(t, err, "Shouldn't attempt to contact any machines")
}

func TestMonitorCantSSH(t *testing.T) {
	t.Parallel()
	// TODO: Mock out pinger to allow this.
	if os.Getenv("CI") != "" {
		t.Skip("Cannot ping on CI")
	}

	db := database.InMemory()

	logger := slog.Default()

	sshConfig := sshConfig(t)

	currentTime := time.Now()
	db.UpdateLastSeen("google.com", currentTime.Add(-48*time.Hour).Unix()) // This machine should trigger actions

	cutoffTime := currentTime.Add(-24 * time.Hour)
	err := monitor(db, cutoffTime, logger, sshConfig)

	if err == nil {
		t.Fatal("Expected an error")
	}

}
