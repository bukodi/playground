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
	pqConn, pqErr := ssh.Dial("tcp", target, clientCfg)
	if pqErr != nil {
		kexInitErr := &ssh.AlgorithmNegotiationError{}
		if errors.As(pqErr, &kexInitErr) {
			result.IsPQKexSupported = "KexInitError"
		} else {
			// Other error
			return result, pqErr
		}
	} else {
		// PQConn successfully established
		result.IsPQKexSupported = true
		pqConn.Close()
	}

	// Connect to the SSH server using non PQ key exchange protocol
	clientCfg.Config.KeyExchanges = []string{"diffie-hellman-group-exchange-sha256"}
	nonPqConn, nonPqErr := ssh.Dial("tcp", target, clientCfg)
	if nonPqErr != nil {
		kexInitErr := &ssh.AlgorithmNegotiationError{}
		if errors.As(nonPqErr, &kexInitErr) {
			result.IsPQKexSupported = "KexInitError"
		} else {
			// Other error
			return result, nonPqErr
		}
	} else {
		// PQConn successfully established
		result.IsNonPQKexSupported = true
		nonPqConn.Close()
	}

	return result, nil

}

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	slog.Info("host key callback", "hostname", hostname, "remote", remote, "key", key)
	return nil
}
