package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func scanTLSPort(host string, port int, timeout time.Duration) (result ScanResult) {
	result.Address = host
	result.Port = port
	now := time.Now()
	defer func() {
		result.TestDuration = time.Since(now)
		if r := recover(); r != nil {
			result = ScanResult{
				Address: host,
				Port:    port,
				Error:   fmt.Sprintf("panic: %v", r),
			}
		}
	}()

	// Create a Dialer with the specified timeout
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// Create a TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set to true for testing purposes only

	}

	// Establish a TLS connection
	conn, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%d", host, port), tlsConfig)
	if err != nil {
		if isNetworkError(err) {
			result.PortType = NoConn
		} else {
			result.PortType = Other
			result.Error = err.Error()
		}
		return
	}
	tlsState := conn.ConnectionState()
	if !tlsState.HandshakeComplete {
		result.PortType = Other
		result.Error = "TLS handshake failed"
		return
	}
	defer conn.Close()
	result.PortType = TLS
	result.TLSVersion = tls.VersionName(tlsState.Version)
	result.CipherSuite = tls.CipherSuiteName(tlsState.CipherSuite)
	result.ServerCertKeyAlgo = tlsState.PeerCertificates[0].PublicKeyAlgorithm.String()
	if curveID := tlsState.CurveID; curveID != 0 {
		result.TLSCurveName = curveName(curveID)
		if curveID == tls.X25519MLKEM768 {
			result.IsPQCCurve = true
		}
	}
	return
}

func curveName(curveID tls.CurveID) string {
	switch curveID {
	case tls.CurveP256:
		return "P-256"
	case tls.CurveP384:
		return "P-384"
	case tls.CurveP521:
		return "P-521"
	case tls.X25519:
		return "X25519"
	case tls.X25519MLKEM768:
		return "X25519MLKEM768"
	case 0x6399:
		return "X25519Kyber768Draft00"
	default:
		return fmt.Sprintf("Unknown Curve ID: %d", curveID)
	}
}
