package tunnel_test

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/gpuctl/gpuctl/internal/tunnel"
)

func TestSSHRestart_UnreadableKey(t *testing.T) {
	t.Parallel()

	// Setup
	sshConfig := tunnel.Config{User: "testuser"}

	err := tunnel.RestartSatellite("localhost", sshConfig)

	if !errors.Is(err, tunnel.InvalidConfigError) {
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

	sshConfig := tunnel.Config{
		User:        "dummyUser",
		Signer:      signer,
		KeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	err = tunnel.RestartSatellite("invalid.remote.address", sshConfig)

	var opError *net.OpError
	if errors.As(err, &opError) && opError.Err != nil {
		t.Logf("Expected error occurred: %v", err)
	} else if err != nil {
		t.Errorf("Expected sshRestart to fail due to 'no such host', but it failed with a different error: %v", err)
	} else {
		t.Error("Expected an error due to invalid setup, but sshRestart did not fail as expected")
	}
}
