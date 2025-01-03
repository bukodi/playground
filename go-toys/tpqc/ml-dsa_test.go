package tpqc

import (
	"testing"
)

func TestMLDSA(t *testing.T) {
	t.Run(ML_DSA_44.Name, func(t *testing.T) {
		if err := testMLDSA(t, ML_DSA_44); err != nil {
			t.Fatal(err)
		}
	})
	t.Run(ML_DSA_65.Name, func(t *testing.T) {
		if err := testMLDSA(t, ML_DSA_65); err != nil {
			t.Fatal(err)
		}
	})
	t.Run(ML_DSA_87.Name, func(t *testing.T) {
		if err := testMLDSA(t, ML_DSA_87); err != nil {
			t.Fatal(err)
		}
	})
}

var defaultSeed [32]byte = [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

func testMLDSA(t *testing.T, alg ML_DSA_Alg) error {
	pubKeyPEM, privKeyPem := GenerateKeyPairWithSeed(alg, defaultSeed)
	t.Log(string(pubKeyPEM))
	t.Log(string(privKeyPem))

	return nil
}
