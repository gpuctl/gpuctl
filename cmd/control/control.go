package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
		User:        conf.SSH.Username,
		DataDir:     conf.SSH.DataDir,
		RemoteConf:  conf.SSH.RemoteConf,
		Signer:      signer,
		KeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Be secure here.
	}

	authenticator := webapi.AuthenticatorFromConfig(conf)
	wa := webapi.NewServer(db, &authenticator, tunnelConf)
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
		err := database.DownsampleOverTime(conf.Database.DownsampleInterval, db)
		errs <- fmt.Errorf("downsampler: %w", err)
	}()
	go func() {
		monitorInterval := time.Duration(conf.Timeouts.MonitorInterval) * time.Second
		deathTimeOut := time.Duration(conf.Timeouts.DeathTimeout) * time.Second
		err := groundstation.MonitorForDeadMachines(monitorInterval, db, deathTimeOut, log, tunnelConf)
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
