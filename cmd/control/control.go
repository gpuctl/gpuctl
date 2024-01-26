package main

import (
	"flag"
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

	var db database.Database
	var err error

	if *inMemDb && *postgresDb {
		fatal("cannot have both `-inmemdb` and `-postgres`")
	} else if *inMemDb {
		db = database.InMemory()
	} else if *postgresDb {
		dbUrl, pres := os.LookupEnv("GPU_DB_URL")
		if !pres {
			fatal("failed to read environment variable GPU_DB_URL")
		}
		db, err = postgres.New(dbUrl)
		if err != nil {
			fatal("failed to connect to database: " + err.Error())
		}
	} else {
		fatal("must pass one of `-inmemdb` and `-postgres`")
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

func fatal(s string) {
	slog.Error(s)
	os.Exit(1)
}
