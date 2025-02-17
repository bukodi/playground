package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"github.com/digitorus/pdf"
	"github.com/digitorus/pdfsign/revocation"
	"github.com/digitorus/pdfsign/sign"
	"log"
	"os"
	"time"
)

func main() {

	certBytes, _ := os.ReadFile("tpdfsign/testdata/test_signer.cer")
	certPem, _ := pem.Decode(certBytes)
	cert, _ := x509.ParseCertificate(certPem.Bytes)

	keyBytes, _ := os.ReadFile("tpdfsign/testdata/test_signer.p8.pem")
	keyPem, _ := pem.Decode(keyBytes)
	anyKey, _ := x509.ParsePKCS8PrivateKey(keyPem.Bytes)
	privateKey := anyKey.(crypto.Signer)

	err := signFn("tpdfsign/testdata/testfile14.pdf", "tpdfsign/testdata/test_signed.pdf", cert, privateKey)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Signed PDF written to ")
	}
}

func signFn(input string, output string, cert *x509.Certificate, privateKey crypto.Signer) error {
	input_file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer input_file.Close()

	output_file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer output_file.Close()

	finfo, err := input_file.Stat()
	if err != nil {
		return err
	}
	size := finfo.Size()

	rdr, err := pdf.NewReader(input_file, size)
	if err != nil {
		return err
	}

	err = sign.Sign(input_file, output_file, rdr, size, sign.SignData{
		Signature: sign.SignDataSignature{
			Info: sign.SignDataSignatureInfo{
				Name:        "John Doe",
				Location:    "Somewhere on the globe",
				Reason:      "My season for siging this document",
				ContactInfo: "How you like",
				Date:        time.Now().Local(),
			},
			CertType:   sign.CertificationSignature,
			DocMDPPerm: sign.AllowFillingExistingFormFieldsAndSignaturesPerms,
		},
		Signer:          privateKey,    // crypto.Signer
		DigestAlgorithm: crypto.SHA256, // hash algorithm for the digest creation
		Certificate:     cert,          // x509.Certificate
		//CertificateChains: certificate_chains, // x509.Certificate.Verify()
		TSA: sign.TSA{
			URL:      "https://freetsa.org/tsr",
			Username: "",
			Password: "",
		},

		// The follow options are likely to change in a future release
		//
		// cache revocation data when bulk signing
		RevocationData: revocation.InfoArchival{},
		// custom revocation lookup
		RevocationFunction: sign.DefaultEmbedRevocationStatusFunction,
	})
	return err
}
