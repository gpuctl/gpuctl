package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/gpuctl/gpuctl/internal/onboard"
)

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

	err = onboard.Onboard(*user, *remote, signer, ssh.InsecureIgnoreHostKey())

	if err != nil {
		log.Fatal(err)
	}
}
