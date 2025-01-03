package tpqc

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"github.com/cloudflare/circl/sign/schemes"
)

type ML_DSA_Alg struct {
	Name string
	Oid  asn1.ObjectIdentifier
}

// https://csrc.nist.gov/projects/computer-security-objects-register/algorithm-registration

var ML_DSA_44 = ML_DSA_Alg{"ML-DSA-44", asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 3, 17}}
var ML_DSA_65 = ML_DSA_Alg{"ML-DSA-65", asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 3, 18}}
var ML_DSA_87 = ML_DSA_Alg{"ML-DSA-87", asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 3, 19}}

type subjectPublicKeyInfo struct {
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

type mldsaPrivateKey struct {
	Version    int
	Algorithm  pkix.AlgorithmIdentifier
	PrivateKey []byte
}

func GenerateKeyPairWithSeed(alg ML_DSA_Alg, seed [32]byte) (publicKey, privateKey []byte) {
	scheme := schemes.ByName(alg.Name)

	for i := 0; i < len(seed); i++ {
		seed[i] = byte(i)
	}

	pubKey, privKey := scheme.DeriveKey(seed[:])
	_ = privKey

	pubKeyBytes, err := pubKey.MarshalBinary()
	if err != nil {
		panic(err)
	}

	pubKeyASN1Bytes, err := asn1.Marshal(subjectPublicKeyInfo{
		Algorithm: pkix.AlgorithmIdentifier{
			Algorithm: alg.Oid,
		},
		PublicKey: asn1.BitString{
			BitLength: len(pubKeyBytes) * 8,
			Bytes:     pubKeyBytes,
		},
	})
	if err != nil {
		panic(err)
	}

	privKeyASN1Bytes, err := asn1.Marshal(mldsaPrivateKey{
		Version: 0,
		Algorithm: pkix.AlgorithmIdentifier{
			Algorithm: alg.Oid,
		},
		PrivateKey: seed[:],
	})
	if err != nil {
		panic(err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyASN1Bytes,
	})

	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyASN1Bytes,
	})

	return pubKeyPEM, privKeyPEM

}
