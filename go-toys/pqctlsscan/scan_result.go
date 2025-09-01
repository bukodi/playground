package main

import (
	"errors"
	"net"
	"syscall"
	"time"
)

type PortType string

const TLS = PortType("TLS")
const SSH = PortType("SSH")
const Other = PortType("Other")
const NoConn = PortType("NoConn")

type ScanResult struct {
	Address             string        `json:"address"`
	Port                int           `json:"port"`
	PortType            PortType      `json:"portType"`
	IsPQKexSupported    bool          `json:"isPQKexSupported"`
	IsNonPQKexSupported bool          `json:"isNonPQKexSupported"`
	Error               string        `json:"error,omitempty"`
	TestDuration        time.Duration `json:"testDuration"`

	TLSVersion        string `json:"tlsVersion"`
	TLSCurveName      string `json:"ecCurve"`
	CipherSuite       string `json:"cipherSuite"`
	ServerCertKeyAlgo string `json:"serverCertKeyAlgo"`
	IsPQCCurve        bool   `json:"isPQCCurve"`
}

func isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	var netErr *net.OpError
	if !errors.As(err, &netErr) {
		return false
	}
	if errors.Is(netErr.Err, syscall.ECONNREFUSED) ||
		errors.Is(netErr.Err, syscall.ECONNRESET) ||
		errors.Is(netErr.Err, syscall.ETIMEDOUT) ||
		errors.Is(netErr.Err, syscall.EHOSTUNREACH) ||
		errors.Is(netErr.Err, syscall.ECONNABORTED) {
		return true
	}
	if netErr.Timeout() {
		return true
	}
	return false
}
