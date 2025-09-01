package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

var pqKexAlogs []string = []string{ssh.KeyExchangeMLKEM768X25519}
var nonPqKeyAlogs = []string{
	ssh.KeyExchangeCurve25519,
	ssh.KeyExchangeECDHP256,
	ssh.KeyExchangeECDHP384,
	ssh.KeyExchangeECDHP521,
	ssh.KeyExchangeDH14SHA256,
	ssh.InsecureKeyExchangeDH14SHA1,
}
var allKexAlgos = append(pqKexAlogs, nonPqKeyAlogs...)

func checkSSHPort(host string, port int, timeout time.Duration) (pqcKexCompleted bool, nonPqcKexCompleted bool, err error) {
	nonPqcKexCompleted, err = checkSSHPortOnce(host, port, false, true, timeout)
	if err != nil {
		return false, false, err
	}
	pqcKexCompleted, err = checkSSHPortOnce(host, port, true, false, timeout)
	if err != nil {
		return false, false, err
	}
	return pqcKexCompleted, nonPqcKexCompleted, nil
}

func checkSSHPortOnce(host string, port int, allowPQCKex bool, allowNonPQCKex bool, timeout time.Duration) (kexCompleted bool, err error) {

	// Create an SSH configuration
	hostKeyCallbackCalled := false
	clientCfg := &ssh.ClientConfig{
		Config: ssh.Config{},
		User:   "cica",
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			hostKeyCallbackCalled = true
			return nil
		},
	}
	if allowPQCKex && !allowNonPQCKex {
		clientCfg.Config.KeyExchanges = pqKexAlogs
	} else if !allowPQCKex && allowNonPQCKex {
		clientCfg.Config.KeyExchanges = nonPqKeyAlogs
	} else if allowPQCKex && allowNonPQCKex {
		clientCfg.Config.KeyExchanges = allKexAlgos
	} else {
		return false, errors.New("invalid argument combination")
	}
	clientCfg.SetDefaults()

	// Establish a TLS connection
	target := fmt.Sprintf("%s:%d", host, port)

	d := net.Dialer{Timeout: timeout}
	rawConn, err := d.Dial("tcp", target)
	if err != nil {
		return false, err
	}
	// Optional: cap total time for SSH handshake as well
	_ = rawConn.SetDeadline(time.Now().Add(timeout))

	sshConn, _, _, err := ssh.NewClientConn(rawConn, target, clientCfg)
	if err != nil {
		rawConn.Close()
		kexInitErr := &ssh.AlgorithmNegotiationError{}
		if errors.As(err, &kexInitErr) {
			return false, nil
		} else if hostKeyCallbackCalled {
			// We got a host key callback, but the authentication failed.
			return true, nil
		} else {
			return false, err
		}
	}
	defer sshConn.Close()
	return true, nil
}

func hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	slog.Info("host key callback", "hostname", hostname, "remote", remote, "key", key)
	return nil
}
