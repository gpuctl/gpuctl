package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"errors"
	"net/http"
	"time"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/gpustats"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

var (
	errSuspectedServerMissingInfo = errors.New("Groundstation could not update it's database with given packet. Likely forgot about this GPU.")
)

func main() {
	log := slog.Default()

	log.Info("Starting satellite")

	host, err := os.Hostname()

	if err != nil {
		log.Error("failed to get hostname", "err", err)
		return
	}

	log.Info("got hostname", "hostname", host)

	satellite_configuration, err := config.GetClientConfiguration("satellite.toml")

	log.Info("got config", "config", satellite_configuration)

	if err != nil {
		log.Error("Failed to get satellite configuration from toml configuration file", "err", err)
		os.Exit(1)
	}

	s := satellite{
		gsAddr: config.GenerateAddress(
			satellite_configuration.Groundstation.Protocol,
			satellite_configuration.Groundstation.Hostname,
			satellite_configuration.Groundstation.Port),
		hostname: host,
	}

	hndlr := setGPUHandler(satellite_configuration.Satellite.FakeGPU)

	// Send initial infopacket of GPUInfo
	log.Info("Sending initial GPU context")
	err = s.sendGPUInfo(hndlr)
	if err != nil {
		log.Error("Failed to send GPU context", "err", err)
	}

	go func() {
		for {
			log.Info("Sending heartbeat")
			err := s.sendHeartBeat()

			if err != nil {
				log.Error("failed to send heartbeat", "err", err)
			}

			time.Sleep(time.Duration(satellite_configuration.Satellite.HeartbeatInterval))
		}
	}()

	// backlog, _ := recoverState(satellite_configuration.Satellite.Cache)

	// for stat := range backlog {
	// 	err := s.sendGPUStatus(backlog[stat])

	// 	if err != nil {
	// 		log.Error("Failed to send backlogged GPU stat message", "err", err)
	// 	}
	// }

	backlog := make([][]uplink.GPUStatSample, 0)

	collectGPUStatTicker := time.NewTicker(time.Duration(satellite_configuration.Satellite.DataInterval))
	publishGPUStatTicker := time.NewTicker(time.Duration(satellite_configuration.Satellite.DataInterval))

	// Go has no API for a ticker with an instantaneous first tick (see
	// https://github.com/golang/go/issues/17601) so we have to use a clunky
	// work-around
	publishGPUStats := func() {
		log.Info("Sending status")

		err = s.sendGPUStatus(processStats(backlog))

		if err == errSuspectedServerMissingInfo {
			// Send context again to try and refresh server
			log.Info("Server could not handle our sample submission. Resubmitting GPU contextual information.")
			err = s.sendGPUInfo(hndlr)
			if err != nil {
				// Just quit, server is either broken or refuses to acknowledge us
				log.Error("Server did not accept our resubmittal of contextual information. Quitting...")
				os.Exit(-1)
			}
		}

		if err != nil {
			log.Error("Failed to publish current GPU stat message", "err", err)
		}
	}

	collectGPUStats := func() {
		log.Info("Collecting GPU Status")

		stat, err := hndlr.GetGPUStatus()

		if err != nil {
			log.Error("Failed to get GPU stat from stat handler", "err", err)
		}

		backlog = append(backlog, stat)
		saveState(satellite_configuration.Satellite.Cache, backlog)
	}

	collectGPUStats()
	publishGPUStats()

	for {
		select {
		case <-publishGPUStatTicker.C:
			publishGPUStats()
		case <-collectGPUStatTicker.C:
			collectGPUStats()
		}
	}
}

type satellite struct {
	hostname string
	gsAddr   string
}

type flags struct {
	fakeGPUs bool
}

func saveState(filename string, items [][]uplink.GPUStatSample) error {
	data, err := json.Marshal(items)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o644)
}

func recoverState(filename string) ([][]uplink.GPUStatSample, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var state [][]uplink.GPUStatSample

	err = json.Unmarshal(data, &state)

	if err != nil {
		return nil, err
	}

	return state, nil
}

// Dummy function to process list of stats
func processStats(stats [][]uplink.GPUStatSample) []uplink.GPUStatSample {
	return stats[len(stats)-1]
}

func setGPUHandler(isFakeGPUs bool) gpustats.GPUDataSource {
	if isFakeGPUs {
		return gpustats.FakeGPU{}
	} else {
		return gpustats.NvidiaGPUHandler{}
	}
}

func (s *satellite) sendHeartBeat() error {
	resp, err := femto.Post(
		s.gsAddr+uplink.HeartbeatUrl,
		uplink.HeartbeatReq{Hostname: s.hostname},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err // TODO: add errors for each status code
	}
	return nil
}

func (s *satellite) sendGPUInfo(gpuhandler gpustats.GPUDataSource) error {
	info, err := gpuhandler.GetGPUInformation()
	if err != nil {
		return err
	}

	resp, err := femto.Post(
		s.gsAddr+uplink.GPUStatsUrl,
		uplink.GpuStatsUpload{Hostname: s.hostname, GPUInfos: info},
	)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err // TODO: Add errors for each code
	}
	return nil
}

func (s *satellite) sendGPUStatusWithSource(gpuhandler gpustats.GPUDataSource) error {
	stats, err := gpuhandler.GetGPUStatus()

	if err != nil {
		return err
	}
	return s.sendGPUStatus(stats)

}

func (s *satellite) sendGPUStatus(stats []uplink.GPUStatSample) error {
	resp, err := femto.Post(
		s.gsAddr+uplink.GPUStatsUrl,
		uplink.GpuStatsUpload{Hostname: s.hostname, Stats: stats},
	)
	if resp.StatusCode == http.StatusBadRequest {
		// Could not send status, suspected cause is server does not have context.
		return errSuspectedServerMissingInfo

	}
	return err // TODO: Handle errors cases
}
