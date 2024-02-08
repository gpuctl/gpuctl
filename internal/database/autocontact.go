package database

import (
	"log/slog"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func MonitorForDeadMachines(interval int, database Database, timespanForDeath int, l *slog.Logger) error {
	downsampleTicker := time.NewTicker(time.Duration(interval))

	for t := range downsampleTicker.C {
		err := monitor(database, t, timespanForDeath, l)

		if err != nil {
			return err
		}
	}

	return nil
}

func monitor(database Database, t time.Time, timespanForDeath int, l *slog.Logger) error {
	last_seens, err := database.LastSeen()

	if err != nil {
		return err
	}

	for idx := range last_seens {
		seen := last_seens[idx]

		if seen.LastSeen < t.Add(-1*time.Duration(timespanForDeath)*time.Second).Unix() {
			if ping(seen.Hostname, l) {
				err := sshRestart(seen.Hostname, l)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func sshRestart(hostname string, l *slog.Logger) error {
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
