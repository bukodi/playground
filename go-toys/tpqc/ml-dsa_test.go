package tpqc

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/cloudflare/circl/sign/schemes"
	"os"
	"testing"
)

type subjectPublicKeyInfo struct {
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

type mldsaPrivateKey struct {
	Version    int
	Algorithm  pkix.AlgorithmIdentifier
	PrivateKey []byte
}

func TestMLDSA(t *testing.T) {
	t.Run("ML-DSA-44", func(t *testing.T) {
		if err := testMLDSA(t, "ML-DSA-44", 17); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("ML-DSA-65", func(t *testing.T) {
		if err := testMLDSA(t, "ML-DSA-65", 18); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("ML-DSA-87", func(t *testing.T) {
		if err := testMLDSA(t, "ML-DSA-87", 19); err != nil {
			t.Fatal(err)
		}
	})
}

func testMLDSA(t *testing.T, name string, oid int) error {
	scheme := schemes.ByName(name)
	var seed [32]byte // 000102…1e1f

	for i := 0; i < len(seed); i++ {
		seed[i] = byte(i)
	}

	pk, _ := scheme.DeriveKey(seed[:])

	ppk, _ := pk.MarshalBinary()

	// https://csrc.nist.gov/projects/computer-security-objects-register/algorithm-registration
	alg := pkix.AlgorithmIdentifier{
		Algorithm: asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 3, oid},
	}

	apk := subjectPublicKeyInfo{
		Algorithm: alg,
		PublicKey: asn1.BitString{
			BitLength: len(ppk) * 8,
			Bytes:     ppk,
		},
	}

	ask := mldsaPrivateKey{
		Version:    0,
		Algorithm:  alg,
		PrivateKey: seed[:],
	}

	papk, err := asn1.Marshal(apk)
	if err != nil {
		return err
	}

	pask, err := asn1.Marshal(ask)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s.pub", name))
	if err != nil {
		return err
	}
	defer f.Close()

	if err = pem.Encode(f, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: papk,
	}); err != nil {
		return err
	}

	f2, err := os.Create(fmt.Sprintf("%s.priv", name))
	if err != nil {
		return err
	}
	defer f2.Close()

	if err = pem.Encode(f2, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pask,
	}); err != nil {
		return err
	}
	return nil
}
