package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/database/postgres"
	"github.com/gpuctl/gpuctl/internal/groundstation"
	"github.com/gpuctl/gpuctl/internal/webapi"
)

func main() {

	inMemDb := flag.Bool("inmemdb", false, "Use an transient, in-memory database")
	postgresDb := flag.Bool("postgres", false, "Use a posgresql database from GPU_DB_URL")

	flag.Parse()

	log := slog.Default()
	log.Info("Starting control server")

	db, err := initialiseDatabase(inMemDb, postgresDb)

	if err != nil {
		fatal("failed to initialise database: " + err.Error())
	}

	gs := groundstation.NewServer(db)
	wa := webapi.NewServer(db)

	errs := make(chan (error), 1)

	go func() {
		errs <- http.ListenAndServe(":8080", gs)
	}()
	go func() {
		errs <- http.ListenAndServe(":8000", wa)
	}()

	slog.Info("started servers")
	err = <-errs
	slog.Error("got an error", "err", err)
}

func initialiseDatabase(inMemDb *bool, postgresDb *bool) (database.Database, error) {
	switch {
	case *inMemDb && *postgresDb:
		return nil, fmt.Errorf("cannot have both '-inmemdb' and '-postgres'")
	case *inMemDb:
		return database.InMemory(), nil
	case *postgresDb:
		dbUrl, pres := os.LookupEnv("GPU_DB_URL")

		if !pres {
			return nil, fmt.Errorf("failed to read enviroment variable GPU_DB_URL")
		}

		return postgres.New(dbUrl)
	default:
		return nil, fmt.Errorf("must pass in one of '-inmemdb' and '-postgres'")
	}
}

func fatal(s string) {
	slog.Error(s)
	os.Exit(1)
}
