package tlstest

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// See https://www.netmeister.org/blog/tls-hybrid-kex.html

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

	s := httptest.NewUnstartedServer(http.HandlerFunc(okHandler))

	// Configure the server to present the certficate we created
	s.TLS = &tls.Config{
		Certificates:       []tls.Certificate{*serverCert},
		GetConfigForClient: nil,
		VerifyConnection:   verifyTLSConnection,
		MinVersion:         tls.VersionTLS13,
		//MaxVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.X25519MLKEM768},
	}

	tls.NewListener(s.Listener, s.TLS)

	// make a HTTPS request to the server
	s.StartTLS()
	defer s.Close()
	t.Logf("Server started on:\n%s\n", s.URL)

	time.Sleep(60 * time.Second)

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

func verifyTLSConnection(state tls.ConnectionState) error {
	fmt.Printf("Server name: %s\n", state.ServerName)
	fmt.Printf("Peer certificates: %v\n", state.PeerCertificates)
	fmt.Printf("Verified chains: %v\n", state.VerifiedChains)
	fmt.Printf("Cipher suite: %x\n", state.CipherSuite)
	fmt.Printf("TLS version: %x\n", state.Version)
	if curveId, err := getTLSCurveID(&state); err != nil {
		fmt.Printf("TLS curve id error:  %s\n", err.Error())
	} else {
		fmt.Printf("TLS curve id:  %x\n", curveId)
	}

	return nil
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	//curveId, err := getRequestCurveID(r)
	curveId, err := getTLSCurveID(r.TLS)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if curveId != tls.X25519MLKEM768 {
		fmt.Fprintf(w, "Insecure connectiom: %x\n\n", curveId)
	} else {
		fmt.Fprintf(w, "Secure PQC connection\n\n")
	}

	w.Write([]byte("HI!"))
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

// getRequestCurveID returns the curve ID of the request
func getTLSCurveID(tlsState *tls.ConnectionState) (tls.CurveID, error) {
	if tlsState == nil {
		return 0, fmt.Errorf("the request is not a TLS connection")
	}

	// Access the private 'testingOnlyCurveID' field using reflection
	connState := reflect.ValueOf(*tlsState)
	curveIDField := connState.FieldByName("testingOnlyCurveID")

	if !curveIDField.IsValid() {
		return 0, fmt.Errorf("the curve ID field is not found")
	}

	// Convert the reflected value to tls.CurveID
	return tls.CurveID(curveIDField.Uint()), nil
}

func getRequestCurveID(r *http.Request) (tls.CurveID, error) {
	if r.TLS == nil {
		return 0, fmt.Errorf("the request is not a TLS connection")
	}

	// Access the private 'testingOnlyCurveID' field using reflection
	connState := reflect.ValueOf(*r.TLS)
	curveIDField := connState.FieldByName("testingOnlyCurveID")

	if !curveIDField.IsValid() {
		return 0, fmt.Errorf("the curve ID field is not found")
	}

	// Convert the reflected value to tls.CurveID
	return tls.CurveID(curveIDField.Uint()), nil
}
