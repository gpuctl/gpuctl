package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/gpuctl/gpuctl/internal/tunnel"
)

// This is right on DoC CSG machines.
// we'll append the user we're logging in as to this path, so we don't run into
// issues overwriting gpuctl directories created by other users in previous
// deployments
const dataDir = "/data/gpuctl"

func main() {
	user := flag.String("user", "", "The username to SSH as")
	keypath := flag.String("key", "", "The path the the SSH private key to authenticate as")
	remote := flag.String("remote", "", "The machine to SSH into")

	flag.Parse()

	if *user == "" {
		log.Fatalf("Must specify `-user`")
	}
	if *keypath == "" {
		log.Fatalf("Must specify `-key`")
	}
	if *remote == "" {
		log.Fatalf("Must specify `-remote`")
	}

	slog.Info("Running onboarding", "user", *user, "key", *keypath)

	key, err := os.ReadFile(*keypath)
	if err != nil {
		log.Fatalf("Unable to read key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse key file %s: %v", *keypath, err)
	}

	conf := tunnel.Config{
		User:    *user,
		DataDir: dataDir + *user,
		Signer:  signer,
		RemoteConf: config.SatelliteConfiguration{
			Groundstation: config.Groundstation{"https://", "gpuctl.perial.co.uk", 80},
			Satellite: config.Satellite{
				DataInterval:      10 * time.Second,
				HeartbeatInterval: 5 * time.Second,
				FakeGPU:           true,
				Cache:             "/data/gpuctl/cache",
			},
		},
	}

	err = tunnel.Onboard(*remote, conf)

	if err != nil {
		log.Fatal(err)
	}
}
