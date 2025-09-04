package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"
)

func checkTLSPort(host string, port int, timeout time.Duration) (result ScanResult) {
	result.Address = host
	result.Port = port
	result.PortType = Other
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

	tlsState, err := checkTLSPortOnce(host, port, false, true, timeout)
	if err != nil && isNetworkError(err) {
		result.PortType = NoConn
		result.Error = err.Error()
		return
	}
	if err == nil && tlsState != nil {
		result.PortType = TLS
		result.IsNonPQKexSupported = true
		result.ServerCertKeyAlgo = tlsState.PeerCertificates[0].PublicKeyAlgorithm.String()
	}
	tlsState, err = checkTLSPortOnce(host, port, true, false, timeout)
	if err != nil && !isNetworkError(err) {
		result.PortType = Other
		return
	}
	if err == nil {
		result.PortType = TLS
		result.IsPQKexSupported = true
		result.ServerCertKeyAlgo = tlsState.PeerCertificates[0].PublicKeyAlgorithm.String()
	}

	return
}

func selectCurves(allowPQCKex bool, allowNonPQCKex bool) ([]tls.CurveID, error) {
	switch {
	case allowPQCKex && !allowNonPQCKex:
		// Force hybrid only
		return []tls.CurveID{tls.X25519MLKEM768}, nil
	case !allowPQCKex && allowNonPQCKex:
		// Classic only
		return []tls.CurveID{tls.X25519, tls.CurveP256}, nil
	case allowPQCKex && allowNonPQCKex:
		// Offer hybrid first, then classic fallback
		return []tls.CurveID{tls.X25519MLKEM768, tls.X25519, tls.CurveP256}, nil
	default:
		return nil, errors.New("no key exchange groups enabled: set AllowPQCipher and/or AllowNonPQCiphers")
	}
}

func checkTLSPortOnce(host string, port int, allowPQCKex bool, allowNonPQCKex bool, timeout time.Duration) (*tls.ConnectionState, error) {
	// Create a Dialer with the specified timeout
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// Create a TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set to true for testing purposes only

	}
	c, err := selectCurves(allowPQCKex, allowNonPQCKex)
	if err != nil {
		return nil, err
	}
	tlsConfig.CurvePreferences = c

	// Establish a TLS connection
	conn, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%d", host, port), tlsConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	tlsState := conn.ConnectionState()
	if !tlsState.HandshakeComplete {
		return nil, errors.New("TLS handshake failed")
	}
	return &tlsState, nil
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
