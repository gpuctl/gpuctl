package groundstation

import (
	"log/slog"
	"os/exec"
	"time"

	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/tunnel"
)

func MonitorForDeadMachines(interval time.Duration, database database.Database, timespanForDeath time.Duration, l *slog.Logger, s tunnel.Config) error {
	downsampleTicker := time.NewTicker(interval)

	for t := range downsampleTicker.C {
		err := monitor(database, t, timespanForDeath, l, s)

		if err != nil {
			l.Error("Error monitoring for dead machines:", "error", err)
		}
	}

	return nil
}

func monitor(database database.Database, t time.Time, timespanForDeath time.Duration, l *slog.Logger, s tunnel.Config) error {
	lastSeens, err := database.LastSeen()

	if err != nil {
		return err
	}

	for idx := range lastSeens {
		seen := lastSeens[idx]

		if seen.LastSeen < t.Add(-1*timespanForDeath*time.Second).Unix() {
			if ping(seen.Hostname, l) {
				err := tunnel.RestartSatellite(seen.Hostname, s)

				if err != nil {
					return err
				}
			}
		}
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
