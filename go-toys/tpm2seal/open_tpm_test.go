package tpm2seal

import (
	"fmt"
	"github.com/google/go-tpm-tools/simulator"
	"github.com/google/go-tpm/tpm2/transport"
	"github.com/google/go-tpm/tpmutil"
	"io"
)

const useSimulator = false

func openTpmDirect() (tpm transport.TPMCloser, retErr error) {
	tpmPath := "/dev/tpmrm0"
	// Open the TPM
	tpm, err := transport.OpenTPM(tpmPath)
	if err != nil {
		return nil, fmt.Errorf("can't open TPM %q: %v", tpmPath, err)
	}

	if !useSimulator {
		return tpm, nil
	}

	sim, err := simulator.Get()
	if err != nil {
		return nil, err
	}
	return &TPM{
		transport: sim,
	}, nil
}

// TPM represents a connection to a TPM simulator.
type TPM struct {
	transport io.ReadWriteCloser
}

// Send implements the TPM interface.
func (t *TPM) Send(input []byte) ([]byte, error) {
	return tpmutil.RunCommandRaw(t.transport, input)
}

// Close implements the TPM interface.
func (t *TPM) Close() error {
	return t.transport.Close()
}
