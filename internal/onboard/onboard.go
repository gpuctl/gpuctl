package onboard

import (
	"bytes"
	"fmt"

	"github.com/gpuctl/gpuctl/internal/assets"
	"github.com/povsister/scp"
	"golang.org/x/crypto/ssh"
)

var remoteDataDir string = "/data/gpuctl"

func Onboard(
	remoteUser string,
	remoteAddr string,
	signer ssh.Signer,
	hostKeyCallback ssh.HostKeyCallback,
) error {
	config := &ssh.ClientConfig{
		User:            remoteUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", remoteAddr+":22", config)
	if err != nil {
		return err
	}
	defer client.Close()

	err = runCommand(client, "mkdir -p "+remoteDataDir)
	if err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	scpClient, err := scp.NewClientFromExistingSSH(client, &scp.ClientOption{})
	if err != nil {
		return err
	}

	err = scpClient.CopyToRemote(
		bytes.NewReader(assets.SatelliteAmd64Linux),
		remoteDataDir+"/satellite",
		&scp.FileTransferOption{Perm: 0o755},
	)
	if err != nil {
		return err
	}

	err = runCommand(client,
		fmt.Sprintf("nohup %s/satellite > %s/satellite.log 2> %s/satellite.err < /dev/null &",
			remoteDataDir, remoteDataDir, remoteDataDir),
	)
	if err != nil {
		return fmt.Errorf("failed to launch satellite on remote: %w", err)
	}

	return nil
}

func runCommand(client *ssh.Client, command string) error {
	sess, err := client.NewSession()
	if err != nil {
		return err
	}

	return sess.Run(command)
}
