package main

import (
	"errors"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestSSHLocalhost(t *testing.T) {

	srv := SSHServer{
		ListenAddr:      "localhost:2222",
		AllowedUser:     "username",
		AllowedPassword: "Passw0rd",
		AllowPQCKex:     false,
		AllowNonPQCKex:  true,
	}
	err := srv.Start(t.Context())
	if err != nil {
		t.Errorf("failed to start SSH server: %+v", err)
		os.Exit(1)
	} else {
		defer srv.Stop()
	}

	if r, err := scanSSHPortWithErr("localhost", 2222, time.Millisecond*100); err != nil {
		kexInitErr := &ssh.AlgorithmNegotiationError{}
		if errors.As(err, &kexInitErr) {
			t.Logf("KexInitError: %v", kexInitErr)
		} else {
			t.Logf("OtherError: %v", err)
		}
	} else {
		t.Logf("Result: %v", r)
	}
}
