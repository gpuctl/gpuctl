package onboard

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/povsister/scp"

	"github.com/gpuctl/gpuctl/internal/assets"
	"github.com/gpuctl/gpuctl/internal/config"
)

// Onboard with copy over and start the satellite on a remote machine, via SSH.
//
// - hostname must be an amd64 linux system
// - user@hostname must have permissions to ssh, when signed in with signer
// - user must have permissions to dataDir
// - dataDir must be machine-local (IE not on NFS)
// - hostKeyCallback will be used to verify the identity of the remote
func Onboard(
	hostname string,
	user string,
	dataDir string,
	signer ssh.Signer,
	keyCallback ssh.HostKeyCallback,
	satConfig config.SatelliteConfiguration,
) error {
	// -- Connect to Remote --
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: keyCallback,
	}
	client, err := ssh.Dial("tcp", hostname+":22", sshConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	// -- Make the data dir --
	err = RunCommand(client, "mkdir -p "+dataDir)
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
		dataDir+"/satellite",
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
		dataDir+"/satellite.toml",
		&scp.FileTransferOption{},
	)
	if err != nil {
		return err
	}

	// -- Start the satellite --
	err = RunCommand(client,
		fmt.Sprintf("nohup %s/satellite >> %s/satellite.log 2>> %s/satellite.err < /dev/null &",
			dataDir, dataDir, dataDir),
	)
	if err != nil {
		return fmt.Errorf("failed to launch satellite on remote: %w", err)
	}

	return nil
}

func Deboard(
	hostname string,
	user string,
	signer ssh.Signer,
	keyCallback ssh.HostKeyCallback,
) error {
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: keyCallback,
	}
	client, err := ssh.Dial("tcp", hostname+":22", sshConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	err = RunCommand(client, "killall satellite")

	if err != nil {
		return fmt.Errorf("failed to kill satellite on remote: %w", err)
	}

	return nil
}

func RunCommand(client *ssh.Client, command string) error {
	sess, err := client.NewSession()
	if err != nil {
		return err
	}

	return sess.Run(command)
}
