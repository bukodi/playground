package tlstest

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func TestCer6766(t *testing.T) {
	// 1. Load the CSR from the specified path
	csrPath := "/home/lbukodi/Downloads/6766.p10"
	csrBytes, err := os.ReadFile(csrPath)
	if err != nil {
		t.Fatalf("Failed to read CSR file: %v", err)
	}

	// 2. Parse the CSR
	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		t.Fatalf("Failed to parse CSR: %v", err)
	}

	// 3. Extract the public key
	pubKey := csr.PublicKey

	// 4. Convert the public key to PKCS8 format
	pkcs8PubKey, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		t.Fatalf("Failed to marshal public key to PKCS8: %v", err)
	}

	// 5. Calculate the SHA-256 hash of the PKCS8 public key
	hash := sha256.Sum256(pkcs8PubKey)
	hashHex := hex.EncodeToString(hash[:])

	// Print the results
	fmt.Printf("CSR Subject: %s\n", csr.Subject)
	fmt.Printf("Public Key Type: %T\n", pubKey)
	fmt.Printf("PKCS8 Public Key Length: %d bytes\n", len(pkcs8PubKey))
	fmt.Printf("SHA-256 Hash: %s\n", hashHex)
}
