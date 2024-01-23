package remote

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/gpustats"
)

// TODO: Remove this
func buildStatusObject(jsonData []byte) (gpustats.GPUStatusPacket, error) {
	var handler gpustats.GPUStatusPacket

	err := json.Unmarshal(jsonData, &handler)
	if err != nil {
		return gpustats.GPUStatusPacket{}, err
	}

	return handler, nil
}

func handleGPUStatusObject(stat gpustats.GPUStatusPacket) error {
	slog.Info("Received packet", "packet", stat)
	return nil
}

func HandleStatusSubmission(packet gpustats.GPUStatusPacket, req *http.Request, log *slog.Logger) error {
	err := handleGPUStatusObject(packet)

	if err != nil {
		return err
	}

	return nil
}
