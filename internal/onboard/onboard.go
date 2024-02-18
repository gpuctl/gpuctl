package onboard

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/povsister/scp"

	"github.com/gpuctl/gpuctl/internal/assets"
	"github.com/gpuctl/gpuctl/internal/config"
)

var InvalidConfigError = errors.New("onboard: invalid config")

type Config struct {
	// The login to run the satellite on other machines as
	User string
	// The directory to store the satellite binary on remotes as
	DataDir string
	// The configuration to install on the remote.
	RemoteConf config.SatelliteConfiguration

	// SSH Options.
	Signer      ssh.Signer
	KeyCallback ssh.HostKeyCallback
}

// Onboard will copy over and start the satellite on a remote machine, via SSH.
//
// - hostname must be an amd64 linux system
// - conf.User@hostname must have permissions to ssh, when signed in with signer
// - conf.User must have permissions to dataDir
// - conf.DataDir must be machine-local (IE not on NFS)
// - conf.HostKeyCallback will be used to verify the identity of the remote
func Onboard(
	hostname string,
	conf Config,
) error {
	// -- Connect to Remote --
	client, err := sshInto(hostname, conf)
	if err != nil {
		return err
	}
	defer client.Close()

	// -- Make the data dir --
	err = runCommand(client, "mkdir -p "+conf.DataDir)
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
		conf.DataDir+"/satellite",
		&scp.FileTransferOption{Perm: 0o755},
	)
	if err != nil {
		return err
	}

	// -- SCP over the config.toml --
	configToml, err := config.ToToml(conf.RemoteConf)
	if err != nil {
		return err
	}
	err = scpClient.CopyToRemote(
		strings.NewReader(configToml),
		conf.DataDir+"/satellite.toml",
		&scp.FileTransferOption{},
	)
	if err != nil {
		return err
	}

	// -- Start the satellite --
	err = startSatellite(client, conf)
	if err != nil {
		return fmt.Errorf("failed to launch satellite on remote: %w", err)
	}

	return nil
}

func RestartSatellite(hostname string, conf Config) error {
	client, err := sshInto(hostname, conf)
	if err != nil {
		return err
	}

	return startSatellite(client, conf)
}

func startSatellite(client *ssh.Client, conf Config) error {
	dataDir := conf.DataDir
	command := fmt.Sprintf(
		"nohup %s/satellite >> %s/satellite.log 2>> %s/satellite.err < /dev/null &",
		dataDir, dataDir, dataDir,
	)
	return runCommand(client, command)
}

func Deboard(
	hostname string,
	conf Config,
) error {
	client, err := sshInto(hostname, conf)
	if err != nil {
		return err
	}
	defer client.Close()

	err = runCommand(client, "killall satellite")
	if err != nil {
		return fmt.Errorf("failed to kill satellite on remote: %w", err)
	}
	return nil
}

func sshInto(hostname string, conf Config) (*ssh.Client, error) {
	if conf.User == "" {
		return nil, fmt.Errorf("%w: User must be set", InvalidConfigError)
	}
	if conf.Signer == nil {
		return nil, fmt.Errorf("%w: Signer must be set", InvalidConfigError)
	}
	if conf.KeyCallback == nil {
		return nil, fmt.Errorf("%w: KetCallback must be set", InvalidConfigError)
	}

	sshConfig := &ssh.ClientConfig{
		User:            conf.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(conf.Signer)},
		HostKeyCallback: conf.KeyCallback,
	}
	return ssh.Dial("tcp", hostname+":22", sshConfig)
}

func runCommand(client *ssh.Client, command string) error {
	sess, err := client.NewSession()
	if err != nil {
		return err
	}

	return sess.Run(command)
}
