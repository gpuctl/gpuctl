package webapi

import (
	"log/slog"
	"net/http"
)

func (wa *api) HandleOfflineMachineRequest(req *http.Request, log *slog.Logger) ([]string, error) {
	machine_data, err := wa.db.LastSeen()

	if err != nil {
		return nil, err
	}

	var names []string

	for idx := range machine_data {
		names = append(names, machine_data[idx].Hostname)
	}

	return names, nil
}
