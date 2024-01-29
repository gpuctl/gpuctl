package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/gpustats"
	"github.com/gpuctl/gpuctl/internal/groundstation/config"
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
	}

	s := satellite{
		// TODO: Make this configurable
		gsAddr: config.GenerateAddress(satellite_configuration.Groundstation.Hostname, satellite_configuration.Groundstation.Port),
		// we assume hostnames don't change during the program's runtime
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

			// TODO: testing only, should not send packets this frequently?
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			log.Info("Sending status")

			err = s.sendGPUStatus(hndlr)

			if err != nil {
				log.Error("failed to send status", "err", err)
			}

			time.Sleep(2 * time.Second)
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

func (s *satellite) sendGPUStatus(gpuhandler gpustats.GPUDataSource) error {
	stats, err := gpuhandler.GPUStats()

	if err != nil {
		return err
	}

	return femto.Post(
		s.gsAddr+uplink.GPUStatsUrl,
		uplink.StatsPackage{Hostname: s.hostname, Stats: stats},
	)
}
