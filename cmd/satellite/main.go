package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/gpustats"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func main() {

	log := slog.Default()

	host, err := os.Hostname()
	if err != nil {
		log.Error("failed to get hostname", "err", err)
		return
	}
	log.Info("got hostname", "hostname", host)

	s := satellite{
		// TODO: Make this configurable
		gsAddr: "http://localhost:8080",
		// we assume hostnames don't change during the program's runtime
		hostname: host,
	}

	log.Info("Starting satellite")

	hndlr := gpustats.NvidiaGPUHandler{}

	for {
		log.Info("Sending heartbeat")
		err := s.sendHeartBeat()
		if err != nil {
			log.Error("failed to send heartbeat", "err", err)
		}
		time.Sleep(2 * time.Second)
		// TODO: testing only, should not send packets this frequently?
		log.Info("Sending status")
		err = s.sendGPUStatus(hndlr)
		if err != nil {
			log.Error("failed to send status", "err", err)
		}
		time.Sleep(2 * time.Second)
	}

	log.Info("Stopped satellite")
}

type satellite struct {
	hostname string
	gsAddr   string
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
