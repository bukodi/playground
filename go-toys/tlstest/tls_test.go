package tlstest

import (
	"crypto/tls"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"
)

//go:embed testdata
var testdata embed.FS

func loadTLSKeyAndCer(name string) (*tls.Certificate, error) {
	cerBytes, err := testdata.ReadFile("testdata/" + name + ".cer")
	if err != nil {
		return nil, err
	}
	keyBytes, err := testdata.ReadFile("testdata/" + name + ".pkcs8")
	if err != nil {
		return nil, err
	}
	tlsCert, err := tls.X509KeyPair(cerBytes, keyBytes)
	return &tlsCert, err
}

func TestTLSServer(t *testing.T) {

	serverCert, err := loadTLSKeyAndCer("test_tls_server")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	ok := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("HI!")) }
	s := httptest.NewUnstartedServer(http.HandlerFunc(ok))

	// Configure the server to present the certficate we created
	s.TLS = &tls.Config{
		Certificates: []tls.Certificate{*serverCert},
		VerifyConnection: func(state tls.ConnectionState) error {
			state.SignedCertificateTimestamps
		},
		GetConfigForClient: fnGetCfgForClient,
	}

	// make a HTTPS request to the server
	s.StartTLS()
	_, err = http.Get(s.URL)
	s.Close()
	t.Logf("%+v", err)

}
