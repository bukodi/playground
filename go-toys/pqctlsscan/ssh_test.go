package main

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	serverHost = "localhost"
	serverPort = 2201
	timeout    = time.Hour * 100
)

func TestSSHKexAlgos(t *testing.T) {
	testCases := []struct {
		name              string
		srvAllowPQCKex    bool
		srvAllowNonPQCKex bool
		expectedPqcOk     bool
		expectedNonPqcOk  bool
	}{
		{
			name:              "Server_nonPQC_only",
			srvAllowPQCKex:    false,
			srvAllowNonPQCKex: true,
			expectedPqcOk:     false,
			expectedNonPqcOk:  true,
		},
		{
			name:              "Server_PQC_only",
			srvAllowPQCKex:    true,
			srvAllowNonPQCKex: false,
			expectedPqcOk:     true,
			expectedNonPqcOk:  false,
		},
		{
			name:              "Server_Both_PQC_and_nonPQC",
			srvAllowPQCKex:    true,
			srvAllowNonPQCKex: true,
			expectedPqcOk:     true,
			expectedNonPqcOk:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			srv := SSHServer{
				ListenAddr:      fmt.Sprintf("%s:%d", serverHost, serverPort),
				AllowedUser:     "username",
				AllowedPassword: "Passw0rd",
				AllowPQCKex:     tc.srvAllowPQCKex,
				AllowNonPQCKex:  tc.srvAllowNonPQCKex,
			}
			if err := srv.Start(t.Context()); err != nil {
				t.Fatalf("failed to start SSH server: %+v (%+v)", err, srv)
				return
			} else {
				defer srv.Stop()
			}

			pqcKexCompleted, nonPqcKexCompleted, err := checkSSHPort(serverHost, serverPort, timeout)
			if err != nil {
				t.Errorf("failed to connect SSH port: %+v", err)
				return
			} else {
				assert.Equal(t, tc.expectedPqcOk, pqcKexCompleted, "PQCKex completed")
				assert.Equal(t, tc.expectedNonPqcOk, nonPqcKexCompleted, "NonPQCKex completed")
			}
		})
	}
}

func TestSSHTimeout(t *testing.T) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	if err != nil {
		t.Fatalf("listen error: %+v", err)
	}
	defer ln.Close()

	_, _, err = checkSSHPort(serverHost, serverPort, time.Millisecond*100)
	if !isNetworkError(err) {
		t.Errorf("expected network error")
		return
	}
}

func TestSSHNetworkError(t *testing.T) {
	_, _, err := checkSSHPort(serverHost, serverPort+1, time.Millisecond*100)
	if !isNetworkError(err) {
		t.Errorf("expected network error")
		return
	}
}
