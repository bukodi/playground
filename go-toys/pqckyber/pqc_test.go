package pqckyber

import (
	"github.com/symbolicsoft/kyber-k2so"
	"testing"
)

func TestPQCKyber(t *testing.T) {
	privAlice, pubAlice, _ := kyberk2so.KemKeypair768()
	t.Logf("Private key (%d bytes): %x", len(privAlice), privAlice)
	t.Logf("Public key  (%d bytes): %x", len(pubAlice), pubAlice)

	t.Logf("----------")

	ciphertext1, ssForSender1, _ := kyberk2so.KemEncrypt768(pubAlice)
	t.Logf("Shared secret on sender side 1   : %x", ssForSender1)
	t.Logf("Encrypted message 1: %x", ciphertext1)
	ssForRecipient1, _ := kyberk2so.KemDecrypt768(ciphertext1, privAlice)
	t.Logf("Shared secret on recipient side 1: %x", ssForRecipient1)

	t.Logf("----------")

	ciphertext2, ssForSender2, _ := kyberk2so.KemEncrypt768(pubAlice)
	t.Logf("Shared secret on sender side 2   : %x", ssForSender2)
	t.Logf("Encrypted message 2: %x", ciphertext2)
	ssForRecipient2, _ := kyberk2so.KemDecrypt768(ciphertext2, privAlice)
	t.Logf("Shared secret on recipient side 2: %x", ssForRecipient2)
}
