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
	dbUrl := configuration.Database.Url
	_, err = postgres.New(dbUrl)
	if err != nil {
		slog.Error("when opening database", "err", err)
		return
	} else {
		slog.Info("connected to database", "url", dbUrl)
	}

	slog.Info("Starting groundstation API server", "port", configuration.Server.Port)

	err = http.ListenAndServe(config.PortToAddress(configuration.Server.Port), srv)

	slog.Info("Shut down groundstation", "err", err)
}
