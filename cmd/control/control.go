package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/groundstation"
	"github.com/gpuctl/gpuctl/internal/tunnel"
	"github.com/gpuctl/gpuctl/internal/webapi"
)

func main() {
	log := slog.Default()
	log.Info("Starting control server")

	conf, err := config.GetControl("control.toml")
	if err != nil {
		fatal("failed to get config: " + err.Error())
	}

	db, err := initialiseDatabase(conf.Database)
	if err != nil {
		fatal("failed to initialise database: " + err.Error())
	}

	gs := groundstation.NewServer(db)
	gsPort := config.PortToAddress(conf.Server.GSPort)

	var signer ssh.Signer
	var key []byte

	if conf.SSH.KeyPath == "" {
		key64 := os.Getenv("GPU_SSH_KEY")
		key, err = base64.StdEncoding.DecodeString(key64)
		if err != nil {
			fatal("failed to decode base64 key: " + err.Error())
		}
	} else {
		key, err = os.ReadFile(conf.SSH.KeyPath)
		if err != nil {
			fatal("failed to read key file: " + err.Error())
		}
	}

	if len(key) != 0 {
		signer, err = ssh.ParsePrivateKey(key)
		if err != nil {
			fatal(fmt.Sprintf("Unable to parse key: %v", err))
		}
	} else {
		// We want to allow the server to run even without an SSH key,
		// for local development.
		log.Warn("No SSH key given, will not be able to handle onboard requests")
	}

	tunnelConf := tunnel.Config{
		User:            conf.SSH.Username,
		DataDirTemplate: conf.SSH.DataDir,
		RemoteConf:      conf.SSH.RemoteConf,
		Signer:          signer,
		KeyCallback:     ssh.InsecureIgnoreHostKey(), // TODO: Be secure here.
	}

	// calculating aggregate data takes a while, so cache it
	var totalEnergy atomic.Uint64
	data, err := db.AggregateData()
	if err != nil {
		log.Error("Got error calculating initial value of aggregate", "err", err)
		os.Exit(-1)
	}
	totalEnergy.Store(data.TotalEnergy)

	const cache_duration = time.Hour
	go func() {
		for range time.NewTicker(cache_duration).C {
			data, err := db.AggregateData()
			if err != nil {
				log.Error("Got error calculating new cache value of aggregate", "err", err)
			}
			totalEnergy.Store(data.TotalEnergy)
		}
	}()

	authenticator := webapi.AuthenticatorFromConfig(conf)
	wa := webapi.NewServer(db, &authenticator, tunnelConf, &totalEnergy)
	waPort := config.PortToAddress(conf.Server.WAPort)

	errs := make(chan (error), 1)

	go func() {
		err := http.ListenAndServe(gsPort, gs)
		errs <- fmt.Errorf("groundstation: %w", err)
	}()
	go func() {
		err := http.ListenAndServe(waPort, wa)
		errs <- fmt.Errorf("webapi: %w", err)
	}()
	go func() {
		// Serve the default mux for pprof debug.
		http.ListenAndServe(":6060", nil)
	}()
	go func() {
		err := database.DownsampleOverTime(conf.Database.DownsampleInterval, db)
		errs <- fmt.Errorf("downsampler: %w", err)
	}()
	go func() {
		err := groundstation.MonitorForDeadMachines(db, conf.Timeouts, log.With(), tunnelConf)
		errs <- fmt.Errorf("dead machine monitor: %w", err)
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
		return database.Postgres(conf.PostgresUrl)
	default:
		return nil, fmt.Errorf("must set one of 'inmemory' or 'postgres'")
	}
}
func fatal(s string) {
	slog.Error(s)
	os.Exit(1)
}
