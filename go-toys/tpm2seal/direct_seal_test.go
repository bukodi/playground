package tpm2seal

import (
	"github.com/google/go-tpm/tpm2"
	"testing"
)

func TestDirectTpmSealUnseal(t *testing.T) {

	thetpm, err := openTpmDirect()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer thetpm.Close()

	// Create the SRK
	// Put a password on the SRK to test more of the flows.

	srkAuth := []byte("password")
	createSRKCmd := tpm2.CreatePrimary{
		PrimaryHandle: tpm2.TPMRHOwner,
		InSensitive: tpm2.TPM2BSensitiveCreate{
			Sensitive: &tpm2.TPMSSensitiveCreate{
				UserAuth: tpm2.TPM2BAuth{
					Buffer: srkAuth,
				},
			},
		},
		InPublic: tpm2.New2B(tpm2.ECCSRKTemplate),
	}
	createSRKRsp, err := createSRKCmd.Execute(thetpm)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("SRK name: %x", createSRKRsp.Name)
	defer func() {
		// Flush the SRK
		flushSRKCmd := tpm2.FlushContext{FlushHandle: createSRKRsp.ObjectHandle}
		if _, err := flushSRKCmd.Execute(thetpm); err != nil {
			t.Errorf("%v", err)
		}
	}()

}
