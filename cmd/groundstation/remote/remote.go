package remote

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/status"
)

// TODO: Remove this
func buildStatusObject(jsonData []byte) (status.GPUStatusPacket, error) {
	var handler status.GPUStatusPacket

	err := json.Unmarshal(jsonData, &handler)
	if err != nil {
		return status.GPUStatusPacket{}, err
	}

	return handler, nil
}

func handleGPUStatusObject(stat status.GPUStatusPacket) error {
	slog.Info("Received packet", "packet", stat)
	return nil
}

func HandleStatusSubmission(packet status.GPUStatusPacket, req *http.Request, log *slog.Logger) error {
	err := handleGPUStatusObject(packet)

	if err != nil {
		return err
	}

	return nil
}
