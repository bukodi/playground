package tlstest

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"io"
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
		//VerifyConnection: func(state tls.ConnectionState) error {
		//	state.SignedCertificateTimestamps
		//},
		//GetConfigForClient: fnGetCfgForClient,
	}

	// make a HTTPS request to the server
	s.StartTLS()
	defer s.Close()

	_, err = http.Get(s.URL)

	clientCert, err := loadTLSKeyAndCer("test_tls_client")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	x509SrvCert, err := x509.ParseCertificate(serverCert.Certificate[0])
	if err != nil {
		t.Fatalf("%+v", err)
	}

	body, err := httpsClientGET(s.URL, clientCert, x509SrvCert)
	if err != nil {
		t.Fatalf("%+v", err)
	} else {
		t.Logf("Body: %s", string(body))
	}

}

func httpsClientGET(url string, clientCert *tls.Certificate, serverCAs ...*x509.Certificate) ([]byte, error) {
	certPool := x509.NewCertPool()
	for _, serverCA := range serverCAs {
		certPool.AddCert(serverCA)
	}
	tlsConfig := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{*clientCert},
	}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return nil, fmt.Errorf("HTTP status: %d, %s", resp.StatusCode, resp.Status)
		}
	}
	msg, err := io.ReadAll(resp.Body)
	return msg, err
}
