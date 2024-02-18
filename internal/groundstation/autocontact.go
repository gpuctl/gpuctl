package groundstation

import (
	"errors"
	"log/slog"
	"os/exec"
	"time"

	"github.com/gpuctl/gpuctl/internal/database"
	"golang.org/x/crypto/ssh"
)

var (
	InvalidSignerError = errors.New("signer is nil")
)

type SSHConfig struct {
	User    string
	Signer  ssh.Signer
	BinPath string
}

func MonitorForDeadMachines(interval time.Duration, database database.Database, timespanForDeath time.Duration, l *slog.Logger, s SSHConfig) error {
	downsampleTicker := time.NewTicker(interval)

	for t := range downsampleTicker.C {
		err := monitor(database, t, timespanForDeath, l, s)

		if err != nil {
			l.Error("Error monitoring for dead machines:", "error", err)
		}
	}

	return nil
}

func monitor(database database.Database, t time.Time, timespanForDeath time.Duration, l *slog.Logger, s SSHConfig) error {
	lastSeens, err := database.LastSeen()

	if err != nil {
		return err
	}

	for idx := range lastSeens {
		seen := lastSeens[idx]

		if seen.LastSeen < t.Add(-1*timespanForDeath*time.Second).Unix() {
			if ping(seen.Hostname, l) {
				err := sshRestart(seen.Hostname, l, s)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func sshRestart(remote string, l *slog.Logger, s SSHConfig) error {
	signer := s.Signer
	user := s.User

	if signer == nil {
		return InvalidSignerError
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		// FIXME: Maybe be more secure??
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", remote+":22", config)
	if err != nil {
		l.Error("Failed to connect to %s: %v", remote, err)
		return err
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		l.Error("Failed to create session: %v", err)
		return err
	}
	defer sess.Close()

	_, err = sess.Output("nohup ./data/gpuctl/satellite")

	if err != nil {
		l.Error("Failed to run command on remote: %s", err)
		return err
	}

	return nil
}

func ping(hostname string, l *slog.Logger) bool {
	cmd := exec.Command("ping", "-c", "1", hostname)
	err := cmd.Run()
	if err != nil {
		l.Debug("Error executing ping command:", err)
		return false
	}

	return cmd.ProcessState.ExitCode() == 0
}
