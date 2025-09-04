package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTLSKexAlgos(t *testing.T) {
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
			srv := TLSServer{
				ListenAddr:        fmt.Sprintf("%s:%d", serverHost, serverPort),
				AllowPQCipher:     tc.srvAllowPQCKex,
				AllowNonPQCiphers: tc.srvAllowNonPQCKex,
				MutualTLSRequired: true,
			}
			if err := srv.Start(t.Context()); err != nil {
				t.Fatalf("failed to start TLS server: %+v (%+v)", err, srv)
				return
			} else {
				defer srv.Stop()
			}

			nonPQState, nonPQErr := checkTLSPortOnce(serverHost, serverPort, false, true, timeout)
			if nonPQErr == nil {
				assert.Equal(t, tc.expectedNonPqcOk, nonPQState != nil, "NonPQKex completed")
			}
			pqState, pqErr := checkTLSPortOnce(serverHost, serverPort, true, false, timeout)
			if pqErr == nil {
				assert.Equal(t, tc.expectedPqcOk, pqState != nil, "PQKex completed")
			}
			if nonPQErr != nil && pqErr != nil {
				t.Fatalf("unexpected error. nonPQErr: %+v, pqErr: %+v", nonPQErr, pqErr)
				assert.Equal(t, nonPQErr.Error(), pqErr.Error(), "Error")
			}
		})
	}
}

func TestTLSTimeout(t *testing.T) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	if err != nil {
		t.Fatalf("listen error: %+v", err)
	}
	defer ln.Close()

	r := checkTLSPort(serverHost, serverPort, time.Millisecond*100)
	assert.Equal(t, NoConn, r.PortType)
}

func TestTLSNetworkError(t *testing.T) {
	r := checkTLSPort(serverHost, serverPort+1, time.Millisecond*100)
	assert.Equal(t, NoConn, r.PortType)
}

func TestNonTLSPort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("world"))
	}))
	defer ts.Close()

	tcpAddr := ts.Listener.Addr().(*net.TCPAddr)

	r := checkTLSPort(tcpAddr.IP.String(), tcpAddr.Port, time.Millisecond*100)
	assert.Equal(t, Other, r.PortType)
}
