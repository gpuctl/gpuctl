package groundstation

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"crypto/rand"
	"crypto/rsa"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/uplink"
	"golang.org/x/crypto/ssh"
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

func (edb *ErrorDB) LatestData() ([]uplink.GpuStatsUpload, error) {
	return nil, nil
}

func (edb *ErrorDB) LastSeen() ([]uplink.WorkstationSeen, error) {
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

func (edb *ErrorDB) Drop() error {
	return nil
}

func TestSSHRestart_UnreadableKey(t *testing.T) {
	// Setup
	logger := slog.Default()
	sshConfig := SSHConfig{
		User: "testuser",
	}

	err := sshRestart("localhost", logger, sshConfig)

	if errors.Is(err, InvalidSignerError) {
		t.Logf("Expected error occurred: %v", err)
	} else if err != nil {
		t.Errorf("Expected sshRestart to fail due to unreadable key, but it failed with a different error: %v", err)
	} else {
		t.Error("Expected an error due to unreadable key, but sshRestart did not fail as expected")
	}
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
	if testing.Short() {
		t.Skip("Skipping ping test in short mode")
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	success := ping("google.com", logger)

	if !success {
		t.Error("Failed to ping google")
	}
}

func TestMonitorWithErrorDB(t *testing.T) {
	db := &ErrorDB{}
	logger := slog.Default()

	sshConfig := SSHConfig{
		User: "testuser",
	}

	currentTime := time.Now()
	timespanForDeath := time.Duration(24*60*60) * time.Second // 24 hours

	err := monitor(db, currentTime, timespanForDeath, logger, sshConfig)
	if err == nil {
		t.Fatal("Expected monitor to return an error due to ErrorDB.LastSeen, but it did not")
	}
}

func TestMonitor(t *testing.T) {
	db := database.InMemory()
	logger := slog.Default()

	sshConfig := SSHConfig{
		User: "testuser",
	}

	currentTime := time.Now()
	db.UpdateLastSeen("machineRecent", currentTime.Unix())                 // This machine should not trigger any action
	db.UpdateLastSeen("machineOld", currentTime.Add(-48*time.Hour).Unix()) // This machine should trigger actions

	timespanForDeath := time.Duration(24*60*60) * time.Second // 24 hours

	err := monitor(db, currentTime, timespanForDeath, logger, sshConfig)
	if err != nil {
		t.Fatalf("Monitor encountered an error: %v", err)
	}
}

func TestSSHRestart_ValidKey(t *testing.T) {
	logger := slog.Default()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA private key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to create SSH signer: %v", err)
	}

	sshConfig := SSHConfig{
		User:   "dummyUser",
		Signer: signer,
	}

	err = sshRestart("invalid.remote.address:22", logger, sshConfig)

	var opError *net.OpError
	if errors.As(err, &opError) && opError.Err != nil {
		t.Logf("Expected error occurred: %v", err)
	} else if err != nil {
		t.Errorf("Expected sshRestart to fail due to 'no such host', but it failed with a different error: %v", err)
	} else {
		t.Error("Expected an error due to invalid setup, but sshRestart did not fail as expected")
	}
}
