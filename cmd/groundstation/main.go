package main

import (
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/cmd/groundstation/config"
	"github.com/gpuctl/gpuctl/internal/database/postgres"
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

	// open database connection
	db, err := postgres.New("postgresql://gpuctl@localhost/gpuctl-tests-db")
	if err != nil {
		slog.Error("opening database:", err)
		return
	}
	_ = db

	slog.Info("Stating groundstation API server", "port", configuration.Server.Port)

	err = http.ListenAndServe(config.PortToAddress(configuration.Server.Port), srv)

	slog.Info("Shut down groundstation", "err", err)
}
