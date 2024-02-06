package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"golang.org/x/crypto/ssh"
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

	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		// FIXME: Maybe be more secure??
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", *remote+":22", config)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", *remote, err)
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer sess.Close()

	out, err := sess.Output("nvidia-smi")
	if err != nil {
		log.Fatalf("Failed to run command on remote: %s", err)
	}

	fmt.Println(string(out))
}
