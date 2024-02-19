package onboard_test

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/gpuctl/gpuctl/internal/onboard"
)

func TestSSHRestart_UnreadableKey(t *testing.T) {
	t.Parallel()

	// Setup
	sshConfig := onboard.Config{User: "testuser"}

	err := onboard.RestartSatellite("localhost", sshConfig)

	if !errors.Is(err, onboard.InvalidConfigError) {
		t.Errorf("expected InvalidConfigError, but got %v", err)
	}
}

func TestSSHRestart_ValidKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA private key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to create SSH signer: %v", err)
	}

	sshConfig := onboard.Config{
		User:        "dummyUser",
		Signer:      signer,
		KeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	err = onboard.RestartSatellite("invalid.remote.address", sshConfig)

	var opError *net.OpError
	if errors.As(err, &opError) && opError.Err != nil {
		t.Logf("Expected error occurred: %v", err)
	} else if err != nil {
		t.Errorf("Expected sshRestart to fail due to 'no such host', but it failed with a different error: %v", err)
	} else {
		t.Error("Expected an error due to invalid setup, but sshRestart did not fail as expected")
	}
}
