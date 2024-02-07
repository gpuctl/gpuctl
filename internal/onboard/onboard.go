package onboard

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gpuctl/gpuctl/internal/assets"
	"github.com/gpuctl/gpuctl/internal/config"
	"github.com/povsister/scp"
	"golang.org/x/crypto/ssh"
)

var remoteDataDir string = "/data/gpuctl"

func Onboard(
	remoteUser string,
	remoteAddr string,
	signer ssh.Signer,
	hostKeyCallback ssh.HostKeyCallback,
	satConfig config.SatelliteConfiguration,
) error {
	// -- Connect to Remote --
	sshConfig := &ssh.ClientConfig{
		User:            remoteUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
	}
	client, err := ssh.Dial("tcp", remoteAddr+":22", sshConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	// -- Make the data dir --
	err = runCommand(client, "mkdir -p "+remoteDataDir)
	if err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	// -- SCP over the satellite binary --
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

	// -- SCP over the config.toml --
	configToml, err := config.ToToml(satConfig)
	if err != nil {
		return err
	}
	err = scpClient.CopyToRemote(
		strings.NewReader(configToml),
		remoteDataDir+"/satellite.toml",
		&scp.FileTransferOption{},
	)
	if err != nil {
		return err
	}

	// -- Start the satellite --
	err = runCommand(client,
		fmt.Sprintf("nohup %s/satellite >> %s/satellite.log 2>> %s/satellite.err < /dev/null &",
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
