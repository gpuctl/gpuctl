package database

import (
	"log/slog"
	"testing"
	"time"
)

func TestSSHRestart(t *testing.T) {
	// Setup
	logger := slog.Default()
	sshConfig := SSHConfig{
		User:    "testuser",
		Keypath: "/path/to/your/test/private/key",
	}

	err := sshRestart("invalid.remote.address", logger, sshConfig)
	if err == nil {
		t.Error("Expected sshRestart to fail due to invalid setup, but it did not")
	} else {
		t.Log("SSHRestart attempted as expected:", err)
	}

}

func TestPing(t *testing.T) {
	logger := slog.Default()

	success := ping("192.0.2.1", logger)

	if success {
		t.Error("Expected ping to fail for a non-routable IP address, but it reported success")
	} else {
		t.Log("Ping processed as expected")
	}
}

func TestMonitor(t *testing.T) {
	db := InMemory()
	logger := slog.Default()

	sshConfig := SSHConfig{
		User:    "testuser",
		Keypath: "/dummy/path",
	}

	currentTime := time.Now()
	db.UpdateLastSeen("machineRecent", currentTime.Unix())                 // This machine should not trigger any action
	db.UpdateLastSeen("machineOld", currentTime.Add(-48*time.Hour).Unix()) // This machine should trigger actions

	timespanForDeath := 24 * 60 * 60 // 24 hours

	err := monitor(db, currentTime, timespanForDeath, logger, sshConfig)
	if err != nil {
		t.Fatalf("Monitor encountered an error: %v", err)
	}

}
