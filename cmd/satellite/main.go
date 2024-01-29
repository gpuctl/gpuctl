package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/gpustats"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func main() {
	flags := parseProgramFlags()

	log := slog.Default()

	log.Info("Starting satellite", "fakegpu", flags.fakeGPUs)

	host, err := os.Hostname()

	if err != nil {
		log.Error("failed to get hostname", "err", err)
		return
	}

	log.Info("got hostname", "hostname", host)

	satellite_configuration, err := config.GetClientConfiguration("config.toml")

	if err != nil {
		log.Error("Failed to get satellite configuration from toml configuration file", "err", err)
		os.Exit(-1)
	}

	s := satellite{
		gsAddr:   config.GenerateAddress(satellite_configuration.Groundstation.Hostname, satellite_configuration.Groundstation.Port),
		hostname: host,
	}

	hndlr := setGPUHandler(flags.fakeGPUs)

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

	go func() {
		backlog, _ := recoverState(satellite_configuration.Satellite.Cache)

		for stat := range backlog {
			err := s.sendGPUStatus(backlog[stat])

			if err != nil {
				log.Error("Failed to send backlogged GPU stat message", "err", err)
			}
		}

		backlog = make([]uplink.GPUStats, 0)

		collectGPUStatTicker := time.NewTicker(time.Duration(satellite_configuration.Satellite.DataInterval) * time.Second)
		publishGPUStatTicker := time.NewTicker(time.Duration(satellite_configuration.Satellite.DataInterval) * time.Second)

		for {
			select {
			case <-publishGPUStatTicker.C:

				log.Info("Sending status")

				err = s.sendGPUStatus(processStats(backlog))

				if err != nil {
					log.Error("Failed to publish current GPU stat message", "err", err)
				}
			case <-collectGPUStatTicker.C:
				log.Info("Collecting GPU Status")

				stat, err := hndlr.GPUStats()

				if err != nil {
					log.Error("Failed to get GPU stat from stat handler", "err", err)
				}

				backlog = append(backlog, stat)
				saveState(satellite_configuration.Satellite.Cache, backlog)
			}

		}
	}()

	log.Info("Stopped satellite")
}

type satellite struct {
	hostname string
	gsAddr   string
}

type flags struct {
	fakeGPUs bool
}

func saveState(filename string, items []uplink.GPUStats) error {
	data, err := json.Marshal(items)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o644)
}

func recoverState(filename string) ([]uplink.GPUStats, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var state []uplink.GPUStats

	err = json.Unmarshal(data, &state)

	if err != nil {
		return nil, err
	}

	return state, nil
}

// Dummy function to process list of stats
func processStats(stats []uplink.GPUStats) uplink.GPUStats {
	return stats[len(stats)-1]
}

func parseProgramFlags() flags {
	fakeGpus := flag.Bool("fakegpu", false, "Use fake GPU data")
	flag.Parse()

	return flags{
		fakeGPUs: *fakeGpus,
	}
}

func setGPUHandler(isFakeGPUs bool) gpustats.GPUDataSource {
	if isFakeGPUs {
		return gpustats.FakeGPU{}
	} else {
		return gpustats.NvidiaGPUHandler{}
	}
}

func (s *satellite) sendHeartBeat() error {
	return femto.Post(
		s.gsAddr+uplink.HeartbeatUrl,
		uplink.HeartbeatReq{Hostname: s.hostname},
	)
}

func (s *satellite) sendGPUStatusWithSource(gpuhandler gpustats.GPUDataSource) error {
	stats, err := gpuhandler.GPUStats()

	if err != nil {
		return err
	}

	return s.sendGPUStatus(stats)

}

func (s *satellite) sendGPUStatus(stats uplink.GPUStats) error {
	return femto.Post(
		s.gsAddr+uplink.GPUStatsUrl,
		uplink.StatsPackage{Hostname: s.hostname, Stats: stats},
	)
}
