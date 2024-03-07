package groundstation

import (
	"log/slog"
	"os/exec"
	"time"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/tunnel"
)

func MonitorForDeadMachines(database database.Database, timeouts config.Timeouts, l *slog.Logger, s tunnel.Config) error {
	downsampleTicker := time.NewTicker(timeouts.MonitorInterval())

	for t := range downsampleTicker.C {
		cutoffTime := t.Add(-timeouts.DeathTimeout())

		err := monitor(database, cutoffTime, l, s)

		if err != nil {
			l.Error("Error monitoring for dead machines:", "error", err)
		}
	}

	return nil
}

// attempts to restart all machines that are reachable, but last pinged up before cutoffTime.
func monitor(database database.Database, cutoffTime time.Time, l *slog.Logger, s tunnel.Config) error {
	lastSeens, err := database.LastSeen()

	if err != nil {
		return err
	}

	for _, seen := range lastSeens {
		// FIXME: If the first machine always fails to restart
		seenIsOld := seen.LastSeen.Before(cutoffTime)
		canPing := ping(seen.Hostname, l)

		l.Info("Deciding if I should restart satellite", "hostname", seen.Hostname, "seen-is-old", seenIsOld, "can-ping", canPing)

		if seenIsOld && canPing {
			l.Info("Attempting to restart a satellite", "hostname", seen.Hostname, "last-seen", seen.LastSeen)

			err := tunnel.RestartSatellite(seen.Hostname, s)

			if err != nil {
				l.Error("Failed to restart machine", "hostname", seen.Hostname, "error", err)
				return err
			}
		}
	}

	return nil
}

func ping(hostname string, l *slog.Logger) bool {
	cmd := exec.Command("ping", "-c", "1", hostname)
	err := cmd.Run()
	if err != nil {
		l.Error("Error executing ping command", "error", err)
		return false
	}

	return cmd.ProcessState.ExitCode() == 0
}
