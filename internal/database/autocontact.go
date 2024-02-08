package database

import (
	"log/slog"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type SSHConfig struct {
	User    string
	Keypath string
}

func MonitorForDeadMachines(interval int, database Database, timespanForDeath int, l *slog.Logger, s SSHConfig) error {
	downsampleTicker := time.NewTicker(time.Duration(interval))

	for t := range downsampleTicker.C {
		err := monitor(database, t, timespanForDeath, l, s)

		if err != nil {
			return err
		}
	}

	return nil
}

func monitor(database Database, t time.Time, timespanForDeath int, l *slog.Logger, s SSHConfig) error {
	last_seens, err := database.LastSeen()

	if err != nil {
		return err
	}

	for idx := range last_seens {
		seen := last_seens[idx]

		if seen.LastSeen < t.Add(-1*time.Duration(timespanForDeath)*time.Second).Unix() {
			if ping(seen.Hostname, l) {
				err := sshRestart(seen.Hostname, l, s)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func sshRestart(remote string, l *slog.Logger, s SSHConfig) error {
	user := s.User
	keypath := s.Keypath

	key, err := os.ReadFile(keypath)
	if err != nil {
		l.Error("Unable to read key: %v", err)
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		l.Error("Unable to parse key file %s: %v", keypath, err)
		return err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		// FIXME: Maybe be more secure??
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", remote+":22", config)
	if err != nil {
		l.Error("Failed to connect to %s: %v", remote, err)
		return err
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		l.Error("Failed to create session: %v", err)
		return err
	}
	defer sess.Close()

	_, err = sess.Output("nohup ./data/gpuctl/satellite")

	if err != nil {
		l.Error("Failed to run command on remote: %s", err)
		return err
	}

	return nil
}

func ping(hostname string, l *slog.Logger) bool {
	// Resolve the hostname to an IP address
	ips, err := net.LookupIP(hostname)
	if err != nil {
		l.Debug("Error resolving hostname:", err)
		return false
	}
	var ip net.IP
	for _, i := range ips {
		if i.To4() != nil {
			ip = i
			break
		}
	}
	if ip == nil {
		l.Debug("No IPv4 address found for hostname")
		return false
	}

	// Create a new ICMP message (Echo Request)
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, // Use PID as ID for simplicity
			Data: []byte("HELLO"),
		},
	}
	binMsg, err := msg.Marshal(nil)
	if err != nil {
		l.Debug("Error marshaling message:", err)
		return false
	}

	// Listen for ICMP replies
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		l.Debug("Error listening for ICMP packets:", err)
		return false
	}
	defer conn.Close()

	// Set a deadline for the ICMP reply
	err = conn.SetDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		l.Debug("Error setting deadline:", err)
		return false
	}

	// Send the ICMP Echo request
	_, err = conn.WriteTo(binMsg, &net.IPAddr{IP: ip})
	if err != nil {
		l.Debug("Error sending ICMP message:", err)
		return false
	}

	// Wait for the Echo reply
	reply := make([]byte, 1500)
	_, _, err = conn.ReadFrom(reply)
	if err != nil {
		l.Debug("Error reading ICMP reply:", err)
		return false
	}

	// If we receive a reply before the deadline, the ping was successful
	return true
}
