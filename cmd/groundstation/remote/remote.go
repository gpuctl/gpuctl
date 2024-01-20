package remote

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/status"
)

func buildStatusObject(jsonData []byte) (status.GPUStatusPacket, error) {
	var handler status.GPUStatusPacket

	err := json.Unmarshal(jsonData, &handler)
	if err != nil {
		return status.GPUStatusPacket{}, err
	}

	return handler, nil
}

func handleGPUStatusObject(stat status.GPUStatusPacket) error {
	return nil
}

func HandleStatusSubmission(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(
			writer, "Invalid method for status submission",
			http.StatusBadRequest,
		)

		return
	}

	body, err := io.ReadAll(request.Body)
	defer request.Body.Close()

	if err != nil {
		http.Error(
			writer, "Expected JSON for status submission",
			http.StatusBadRequest,
		)
	}

	packet, err := buildStatusObject(body)

	if err != nil {
		http.Error(
			writer, "JSON deserialisation was not successful",
			http.StatusBadRequest,
		)
	}

	err = handleGPUStatusObject(packet)

	if err != nil {
		http.Error(
			writer, "There was an error while handling the status object",
			http.StatusInternalServerError,
		)
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK: Submission processed successfully"))
}
