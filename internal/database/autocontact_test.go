package database

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

func TestSSHRestart_UnreadableKey(t *testing.T) {
	logger := slog.Default()
	sshConfig := SSHConfig{
		User:    "testuser",
		Keypath: "/path/to/nonexistent/key", // Intentionally wrong to trigger the read error
	}

	err := sshRestart("localhost", logger, sshConfig)
	if err == nil {
		t.Error("Expected an error due to unreadable key, but got none")
	} else {
		t.Log("Received expected error:", err)
	}
}

func TestSSHRestart_UnparsableKey(t *testing.T) {
	logger := slog.Default()
	// Create a temp file with invalid SSH key content
	tmpFile, err := os.CreateTemp("", "invalid-ssh-key-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	_, err = tmpFile.WriteString("not-a-valid-key")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	sshConfig := SSHConfig{
		User:    "testuser",
		Keypath: tmpFile.Name(), // Use the temp file with invalid content
	}

	err = sshRestart("localhost", logger, sshConfig)
	if err == nil {
		t.Error("Expected an error due to unparsable key, but got none")
	} else {
		t.Log("Received expected error:", err)
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

func TestSSHRestart_ValidKey(t *testing.T) {
	logger := slog.Default()

	// Generate an RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA private key: %v", err)
	}

	// Encode the private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	// Create a temporary file for the SSH key
	tmpFile, err := os.CreateTemp("", "ssh-key-*")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after the test

	// Write the PEM-encoded private key to the temporary file
	if _, err := tmpFile.Write(privateKeyPEM); err != nil {
		t.Fatalf("Failed to write private key to temporary file: %v", err)
	}
	tmpFile.Close()

	// Setup SSHConfig to use the temporary key
	sshConfig := SSHConfig{
		User:    "dummyUser",
		Keypath: tmpFile.Name(),
	}

	// The test below assumes SSHRestart will fail due to invalid remote address but will pass key loading step.
	err = sshRestart("invalid.remote.address:22", logger, sshConfig)
	if err == nil {
		t.Error("Expected SSHRestart to fail due to invalid remote address, but it did not")
	} else {
		t.Log("SSHRestart failed as expected:", err)
	}
}
