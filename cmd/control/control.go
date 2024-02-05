package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/database/postgres"
	"github.com/gpuctl/gpuctl/internal/groundstation"
	"github.com/gpuctl/gpuctl/internal/webapi"
)

func main() {
	log := slog.Default()
	log.Info("Starting control server")

	conf, err := config.GetServerConfiguration("control.toml")
	if err != nil {
		fatal("failed to get config: " + err.Error())
	}

	db, err := initialiseDatabase(conf.Database)
	if err != nil {
		fatal("failed to initialise database: " + err.Error())
	}

	gs := groundstation.NewServer(db)
	gs_port := config.PortToAddress(conf.Server.GSPort)
	wa := webapi.NewServer(db)
	wa_port := config.PortToAddress(conf.Server.WAPort)

	errs := make(chan (error), 1)

	go func() {
		errs <- http.ListenAndServe(gs_port, gs)
	}()
	go func() {
		errs <- http.ListenAndServe(wa_port, wa)
	}()

	slog.Info("started servers")
	err = <-errs
	slog.Error("got an error", "err", err)
}

func initialiseDatabase(conf config.Database) (database.Database, error) {
	switch {
	case conf.InMemory && conf.Postgres:
		return nil, fmt.Errorf("cannot have both 'inmemory' and 'postgres' set")
	case conf.InMemory:
		return database.InMemory(), nil
	case conf.Postgres:
		return postgres.New(conf.PostgresUrl)
	default:
		return nil, fmt.Errorf("must set one of 'inmemory' or 'postgres'")
	}
}

func fatal(s string) {
	slog.Error(s)
	os.Exit(1)
}
