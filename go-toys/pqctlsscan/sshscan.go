package main

import (
	"errors"
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

	kexAlgos := []string{ssh.KeyExchangeMLKEM768X25519}

	// Create an SSH configuration
	clientCfg := &ssh.ClientConfig{
		HostKeyCallback: hostKeyCallback,
		Config: ssh.Config{
			KeyExchanges: kexAlgos,
		},
	}
	clientCfg.SetDefaults()

	// Establish a TLS connection
	target := fmt.Sprintf("%s:%d", host, port)

	// Connect to the SSH server using PQ key exchange protocol
	pqConn, err := ssh.Dial("tcp", target, clientCfg)
	if err != nil {
		kexInitErr := &ssh.AlgorithmNegotiationError{}
		if errors.As(err, &kexInitErr) {
			t.Logf("KexInitError: %v", kexInitErr)
		} else {
			t.Logf("OtherError: %v", err)
		}
		return result, err
	}
	defer pqConn.Close()
	return result, nil
}

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	slog.Info("host key callback", "hostname", hostname, "remote", remote, "key", key)
	return nil
}
