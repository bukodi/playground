package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

func scanSSHPortWithErr(host string, port int, timeout time.Duration) (ScanResult, error) {
	result := ScanResult{
		Address: host,
		Port:    port,
	}

	// Create an SSH configuration
	clientCfg := &ssh.ClientConfig{
		User: "username",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
			// or use ssh.PublicKeys(key) for key-based auth
		},
		HostKeyCallback: hostKeyCallback,
	}
	clientCfg.SetDefaults()

	// Establish a TLS connection
	target := fmt.Sprintf("%s:%d", host, port)

	// Connect to the SSH server
	conn, err := ssh.Dial("tcp", target, clientCfg)
	if err != nil {
		return result, err
	}
	defer conn.Close()
	return result, nil
}

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	slog.Info("host key callback", "hostname", hostname, "remote", remote, "key", key)
	return nil
}
