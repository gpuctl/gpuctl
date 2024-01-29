package main

import (
	"flag"
	"os/signal"
	"log/slog"
	"os"
	"time"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/gpustats"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func main() {
	fakeGpus := flag.Bool("fakegpu", false, "Use fake GPU data")

	// don't define flags below here
	flag.Parse()
	// don't access flags above here

	log := slog.Default()

	log.Info("Starting satellite", "fakegpu", *fakeGpus)

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

	var hndlr gpustats.GPUDataSource
	if *fakeGpus {
		hndlr = gpustats.FakeGPU{}
	} else {
		hndlr = gpustats.NvidiaGPUHandler{}
	}

	// Send initial infopacket of GPUInfo
	s.sendGPUInfo(hndlr)

	// Set up neat ctrl-c breaking
	termsig := make(chan os.Signal, 1)
	signal.Notify(termsig, os.Interrupt)

	go func() {
		<-termsig
		log.Info("Stopping satellite...")
		os.Exit(0)
	}()

	// Start sending heartbeats and statuses
	for {
		log.Info("Sending heartbeat")
		err := s.sendHeartBeat()
		if err != nil {
			log.Error("Failed to send heartbeat", "err", err)
		}
		time.Sleep(2 * time.Second)
		// TODO: testing only, should not send packets this frequently?
		log.Info("Sending status")
		err = s.sendGPUStatus(hndlr)
		if err != nil {
			log.Error("Failed to send status", "err", err)
		}
		time.Sleep(2 * time.Second)
	}
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

func (s *satellite) sendGPUInfo(gpuhandler gpustats.GPUDataSource) error {
	info, err := gpuhandler.GetGPUInformation()
	if err != nil {
		return err
	}

	return femto.Post(
		s.gsAddr+uplink.GPUStatsUrl,
		uplink.GpuStatsUpload{Hostname: s.hostname, GPUInfos: info},
	)
}

func (s *satellite) sendGPUStatus(gpuhandler gpustats.GPUDataSource) error {
	stats, err := gpuhandler.GetGPUStatus()

	if err != nil {
		return err
	}

	return femto.Post(
		s.gsAddr+uplink.GPUStatsUrl,
		uplink.GpuStatsUpload{Hostname: s.hostname, Stats: stats},
	)
}
