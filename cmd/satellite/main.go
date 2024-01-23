package main

import (
	"log/slog"
	"time"

	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/status/handlers"
	"github.com/gpuctl/gpuctl/internal/uplink"
)

func main() {

	log := slog.Default()

	s := satellite{
		// TODO: Make this configurable
		gsAddr: "http://localhost:8080",
	}

	log.Info("Starting satellite")

	hndlr := handlers.NvidiaGPUHandler{}

	for i := 0; i < 10; i++ {
		log.Debug("Sending packets")
		err := s.sendHeartBeat()
		if err != nil {
			log.Error("failed to send heartbeat", "err", err)
		}
		// TODO: testing only, should not send packets this frequently?
		err = s.sendGPUStatus(hndlr)
		if err != nil {
			log.Error("failed to send status", "err", err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	log.Info("Stopped satellite")

}

type satellite struct {
	gsAddr string
}

func (s *satellite) sendHeartBeat() error {
	return femto.Post(
		s.gsAddr+uplink.HeartbeatUrl,
		uplink.HeartbeatReq{Time: time.Now()},
	)
}

func (s *satellite) sendGPUStatus(gpuhandler handlers.GPUDataSource) error {
	packet, err := gpuhandler.GetGPUStatus()

	if err != nil {
		return err
	}

	return femto.Post(
		s.gsAddr+uplink.StatusSubmissionUrl,
		packet,
	)
}
