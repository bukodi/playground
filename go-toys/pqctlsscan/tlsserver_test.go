package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"
)

// TLSServer struct (as referenced in your tests)
type TLSServer struct {
	ListenAddr        string
	AllowPQCipher     bool
	AllowNonPQCiphers bool
	MutualTLSRequired bool
	ln                net.Listener
	tlsLn             net.Listener
	serveCh           chan struct{}
}

func (s *TLSServer) Start(ctx context.Context) error {
	// Select TLS 1.3 groups based on flags
	curves, err := s.selectCurves()
	if err != nil {
		return err
	}

	// Create a minimal HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var msg string
		if r.TLS != nil && r.TLS.CurveID == tls.X25519MLKEM768 {
			// We just echo the selected key exchange group via the server's config lookup
			msg = "PQC key exchange\n"
		} else {
			msg = "Non PQ proof key exchange\n"
		}
		_, _ = w.Write([]byte(msg))
	})

	// Generate an ephemeral self-signed cert for quick use
	cert, err := generateSelfSignedCert()
	if err != nil {
		return fmt.Errorf("failed to generate self-signed cert: %w", err)
	}

	tlsCfg := &tls.Config{
		MinVersion:       tls.VersionTLS13,
		MaxVersion:       tls.VersionTLS13,
		Certificates:     []tls.Certificate{cert},
		CurvePreferences: curves,
	}
	if s.MutualTLSRequired {
		// For mTLS in tests, accept any valid client certificate.
		// In production, set ClientCAs and use tls.RequireAndVerifyClientCert instead.
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	s.ln = ln

	server := &http.Server{
		Addr:      s.ListenAddr,
		Handler:   mux,
		TLSConfig: tlsCfg,
	}

	// Close server when context is cancelled
	go func() {
		<-ctx.Done()
		_ = server.Close()
	}()

	s.tlsLn = tls.NewListener(ln, tlsCfg)
	s.serveCh = make(chan struct{})

	go func() {
		defer close(s.serveCh)
		_ = server.Serve(s.tlsLn)
	}()

	return nil
}

func (s *TLSServer) Stop() {
	// Close listeners; Serve will return and serveCh will close
	if s.tlsLn != nil {
		_ = s.tlsLn.Close()
	}
	if s.ln != nil {
		_ = s.ln.Close()
	}
	// Wait for the serving goroutine to exit
	if s.serveCh != nil {
		<-s.serveCh
	}
}

func (s *TLSServer) selectCurves() ([]tls.CurveID, error) {
	switch {
	case s.AllowPQCipher && !s.AllowNonPQCiphers:
		// Force hybrid only
		return []tls.CurveID{tls.X25519MLKEM768}, nil
	case !s.AllowPQCipher && s.AllowNonPQCiphers:
		// Classic only
		return []tls.CurveID{tls.X25519}, nil
	case s.AllowPQCipher && s.AllowNonPQCiphers:
		// Offer hybrid first, then classic fallback
		return []tls.CurveID{tls.X25519MLKEM768, tls.X25519}, nil
	default:
		return nil, errors.New("no key exchange groups enabled: set AllowPQCipher and/or AllowNonPQCiphers")
	}
}

// generateSelfSignedCert creates a short-lived ECDSA P-256 self-signed certificate.
func generateSelfSignedCert() (tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 62))
	if err != nil {
		return tls.Certificate{}, err
	}

	tpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"TLSServer"},
		},
		NotBefore:             time.Now().Add(-1 * time.Minute),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	cert := tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  priv,
	}
	return cert, nil
}
