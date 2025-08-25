package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"reflect"
	"time"
)

func scanTLSPortWithErr(host string, port int, timeout time.Duration) (ScanResult, error) {
	result := ScanResult{
		Address: host,
		Port:    port,
	}
	// Create a Dialer with the specified timeout
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// Create a TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set to true for testing purposes only

	}

	// Establish a TLS connection
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.DialWithDialer(dialer, "tcp", target, tlsConfig)
	if err != nil {
		return result, err
	}
	tlsState := conn.ConnectionState()
	if !tlsState.HandshakeComplete {
		result.Error = "TLS handshake failed"
		return result, nil
	}
	defer conn.Close()
	result.TLSVersion = tls.VersionName(tlsState.Version)
	result.CipherSuite = tls.CipherSuiteName(tlsState.CipherSuite)
	result.ServerCertKeyAlgo = tlsState.PeerCertificates[0].PublicKeyAlgorithm.String()
	//if curveID, err := getTLSCurveID(&tlsState); err == nil {
	if curveID := tlsState.CurveID; curveID != 0 {
		result.CurveName = curveName(curveID)
		if curveID == tls.X25519MLKEM768 {
			result.IsPQCCurve = true
		}
	}

	return result, nil
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

// deprecated: Emulates tlsState.CurveID before go 1.25
func getTLSCurveID(tlsState *tls.ConnectionState) (tls.CurveID, error) {
	if tlsState == nil {
		return 0, fmt.Errorf("the request is not a TLS connection")
	}
	// Access the private 'testingOnlyCurveID' field using reflection
	connState := reflect.ValueOf(*tlsState)
	curveIDField := connState.FieldByName("testingOnlyCurveID")

	if !curveIDField.IsValid() {
		return 0, fmt.Errorf("the curve ID field is not found")
	}

	// Convert the reflected value to tls.CurveID
	return tls.CurveID(curveIDField.Uint()), nil
}
