package webapi

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/femto"
)

func (wa *Api) HandleOfflineMachineRequest(req *http.Request, log *slog.Logger) (*femto.Response[[]string], error) {
	machine_data, err := wa.DB.LastSeen()

	if err != nil {
		return nil, err
	}

	var names []string

	for idx := range machine_data {
		names = append(names, machine_data[idx].Hostname)
	}

	response := femto.Ok(names)
	return &response, nil
}
