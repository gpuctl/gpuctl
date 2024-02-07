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

	sshClient, err := ssh.Dial("tcp", remoteAddr+":22", config)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sess, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	err = sess.Run("mkdir -p " + remoteDataDir)
	if err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	scpCon, err := scp.NewClientFromExistingSSH(sshClient, &scp.ClientOption{})
	if err != nil {
		return err
	}

	err = scpCon.CopyToRemote(
		bytes.NewReader(assets.SatelliteAmd64Linux),
		remoteDataDir+"/satellite",
		&scp.FileTransferOption{Perm: 0o755},
	)
	if err != nil {
		return err
	}

	return nil
}
