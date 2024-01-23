package main

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/cmd/groundstation/config"
	"github.com/gpuctl/gpuctl/internal/groundstation"
)

func main() {
	slog.Info("Starting groundstation")

	configuration, err := config.GetConfiguration("config.toml")

	if err != nil {
		slog.Error("failed to load config, shutting down", "err", err)
		return
	}

	srv := groundstation.NewServer()

	slog.Info("Stating groundstation API server", "port", configuration.Server.Port)

	err = http.ListenAndServe(config.PortToAddress(configuration.Server.Port), srv)

	slog.Info("Shut down groundstation", "err", err)
}
