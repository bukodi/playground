package rfc3161_test

import (
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"github.com/bukodi/playground/rfc3161"
	"os"
	"testing"
	"time"
)

func TestTSA(t *testing.T) {

	msg := []byte("Hello world!")
	digest := sha256.Sum256(msg)

	tsreq, err := rfc3161.NewTimeStampReq(crypto.SHA256, digest[:])
	if err != nil {
		t.Fatalf("%+v", err)
	}
	tsreq.CertReq = true
	c := rfc3161.NewClient("https://freetsa.org/tsr")
	tsrsp, err := c.Do(tsreq)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	signPemData, err := os.ReadFile("testdata/freetsa_rfc3161_sign.crt")
	pemBlock, _ := pem.Decode(signPemData)
	signCer, err := x509.ParseCertificate(pemBlock.Bytes)

	caPemData, err := os.ReadFile("testdata/freetsa_rootca.crt")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	rfc3161.RootCerts = x509.NewCertPool()
	rfc3161.RootCerts.AppendCertsFromPEM(caPemData)

	err = tsrsp.Verify(tsreq, signCer)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	tsInfo, err := tsrsp.GetTSTInfo()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("Generation Time: %s", tsInfo.GenTime.Format(time.RFC3339Nano))
	t.Logf("Local Time: %s", tsInfo.GenTime.Local().Format(time.DateTime))

}
