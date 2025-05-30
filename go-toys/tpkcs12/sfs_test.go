package tpkcs12

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"testing"
)

func TestInput(t *testing.T) {
	inputData, err := hex.DecodeString(inputHexString)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	sha256Bytes := sha256.Sum256(inputData)
	t.Logf("sha256Bytes as hex: %s", hex.EncodeToString(sha256Bytes[:]))
	t.Logf("sha256Bytes as b64: %s", base64.StdEncoding.EncodeToString(sha256Bytes[:]))

	p7Block, rest := pem.Decode([]byte(p7bPEM))
	if len(rest) != 0 {
		t.Fatalf("decode failed: %+v", rest)
	}
	p7Bytes := p7Block.Bytes

	embeddedCertBytes, err := base64.RawURLEncoding.DecodeString(embeddedCertB64)
	if err != nil {
		t.Errorf("%+v", err)
		t.Logf("embeddedCertB64[615:]: %s", embeddedCertB64[615:700])
	} else if bytes.Equal(p7Bytes[60:60+1817], embeddedCertBytes) {
		t.Logf("embeddedCertBytes matches")
	} else {
		t.Logf("embeddedCertBytes does not match")
	}
	embeddedCert, err := x509.ParseCertificate(embeddedCertBytes)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("embeddedCert.Subject.CommonName: %s", embeddedCert.Subject.CommonName)

	signedAttrsBytes, err := base64.RawURLEncoding.DecodeString(signedAttrsB64)
	//signedAttrsBytes = signedAttrsBytes[2:]
	if err != nil {
		t.Errorf("%+v", err)
	} else if bytes.Equal(p7Bytes[2106:2106+107], signedAttrsBytes) {
		t.Logf("signedAttrsBytes matches")
	} else {
		t.Fatalf("signedAttrsBytes does not match")
	}
	/*msgDigestAttrBytes, err := base64.RawURLEncoding.DecodeString(msgDigestAttrB64)
	//signedAttrsBytes = signedAttrsBytes[2:]
	if err != nil {
		t.Errorf("%+v", err)
	} else if bytes.Equal(p7Bytes[2164:2164+49], msgDigestAttrBytes) {
		t.Logf("msgDigestAttrBytes matches")
	} else {
		t.Fatalf("msgDigestAttrBytes do not match")
	}
	msgDigestAttrBytes = msgDigestAttrBytes[4:] */
	signedAttrsBytes[0] = 0x31
	signedAttrsHash := sha256.Sum256(signedAttrsBytes)

	signatureBytes, err := base64.RawURLEncoding.DecodeString(signatureB64)
	signatureBytes = signatureBytes[4:]
	if err != nil {
		t.Errorf("%+v", err)
	} else if bytes.Equal(p7Bytes[2228+4:2228+4+256], signatureBytes) {
		t.Logf("signatureBytes matches")
	} else {
		t.Fatalf("signatureBytes do not match")
	}

	rsaPubKey := embeddedCert.PublicKey.(*rsa.PublicKey)
	//t.Logf("rsaPubKey: %+v", rsaPubKey)
	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, signedAttrsHash[:], signatureBytes)
	//err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, inputData, signatureBytes)
	if err != nil {
		t.Errorf("Verification failed: %+v", err)
	} else {
		t.Fatalf("signature verified")
	}

}

const inputHexString = "86acc6298c708b16cc24d3f09608ddfdfc24c1a57c3ccf938d4d6df0913ffc9c"

const signedAttrsB64 = `oGkwGAYJKoZIhvcNAQkDMQsGCSqGSIb3DQEHATAcBgkqhkiG9w0BCQUxDxcNMTUwMjA2MTAyMjM3WjAvBgkqhkiG9w0BCQQxIgQghqzGKYxwixbMJNPwlgjd_fwkwaV8PM-TjU1t8JE__Jw`

const msgDigestAttrB64 = `MC8GCSqGSIb3DQEJBDEiBCCGrMYpjHCLFswk0_CWCN39_CTBpXw8z5ONTW3wkT_8nA`

const embeddedCertB64 = `MIIHFTCCBf2gAwIBAgIOTQkYvAd2pLRlx7SyaLswDQYJKoZIhvcNAQELBQAwgbUxCzAJBgNVBAYTAkhVMREwDwYDVQQHDAhCdWRhcGVzdDEVMBMGA1UECgwMTmV0TG9jayBLZnQuMTcwNQYDVQQLDC5UYW7DunPDrXR2w6FueWtpYWTDs2sgKENlcnRpZmljYXRpb24gU2VydmljZXMpMUMwQQYDVQQDDDpOZXRMb2NrIEV4cHJlc3N6IEVhdC4gKENsYXNzIEMgTGVnYWwpIFRhbsO6c8OtdHbDoW55a2lhZMOzMB4XDTEzMTIwNTEzMTAzNFoXDTE1MTIwNTEzMTAzNFowgZIxCzAJBgNVBAYTAkhVMREwDwYDVQQHDAhCdWRhcGVzdDEdMBsGA1UECgwUTWFneWFyIFRlbGVrb20gTnlydC4xDjAMBgNVBAsMBVRFU1pUMR0wGwYDVQQDDBRNYWd5YXIgVGVsZWtvbSBOeXJ0LjEiMCAGCSqGSIb3DQEJARYTcGtpLXRlYW1AdGVsZWtvbS5odTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOYLluDTTsyQOqgY6LupOZ1EIdFPV_1rOMWovpqxpFXJBLcAnH5Dvt-RyIWzvTOe_ALzdgtPu9C0KjtB7qpzaHbcsTc7PLmqhwOtHYUacSNBZU2bmK8rBr7r5irIC3-0Q8o-_qMYMqWTY__nXLCo6-aM8-xTIOp5otdZjxBK7jpWeFRJerLkjnXFVHSQIPxj7GQn_DsxI97cRPDgyQJcLgEb4OfXu4srIk4ymrpaIj8YIZ06yBBnhg4a_XC77MVlUIw_rDTMtv2mFotxn6w6Wbt2v9Zg1C1o8ROY5n8WPQFPjxxqgGzBi7ssxhKIPkp2ab6U5h2DnDs3qxSztgIihckCAwEAAaOCA0IwggM-MAwGA1UdEwEB_wQCMAAwDgYDVR0PAQH_BAQDAgbAMB0GA1UdDgQWBBT1CZ_aem6MvjoFOXTiI-1ytxxfoTAfBgNVHSMEGDAWgBRxSW6Bh9EIrfEoBKHzv_9Gb3-12jAdBgNVHSUEFjAUBggrBgEFBQcDBAYIKwYBBQUHAwIwHgYDVR0RBBcwFYETcGtpLXRlYW1AdGVsZWtvbS5odTCCAUQGCCsGAQUFBwEBBIIBNjCCATIwLQYIKwYBBQUHMAGGIWh0dHA6Ly9vY3NwMS5uZXRsb2NrLmh1L2NjbGNhLmNnaTAtBggrBgEFBQcwAYYhaHR0cDovL29jc3AyLm5ldGxvY2suaHUvY2NsY2EuY2dpMC0GCCsGAQUFBzABhiFodHRwOi8vb2NzcDMubmV0bG9jay5odS9jY2xjYS5jZ2kwNQYIKwYBBQUHMAKGKWh0dHA6Ly9haWExLm5ldGxvY2suaHUvaW5kZXguY2dpP2NhPWNjbGNhMDUGCCsGAQUFBzAChilodHRwOi8vYWlhMi5uZXRsb2NrLmh1L2luZGV4LmNnaT9jYT1jY2xjYTA1BggrBgEFBQcwAoYpaHR0cDovL2FpYTMubmV0bG9jay5odS9pbmRleC5jZ2k_Y2E9Y2NsY2EwgaEGA1UdHwSBmTCBljAwoC6gLIYqaHR0cDovL2NybDEubmV0bG9jay5odS9pbmRleC5jZ2k_Y3JsPWNjbGNhMDCgLqAshipodHRwOi8vY3JsMi5uZXRsb2NrLmh1L2luZGV4LmNnaT9jcmw9Y2NsY2EwMKAuoCyGKmh0dHA6Ly9jcmwzLm5ldGxvY2suaHUvaW5kZXguY2dpP2NybD1jY2xjYTCBsgYDVR0gBIGqMIGnMIGkBg0rBgEEAZtjASaJyvBRMIGSMCcGCCsGAQUFBwIBFhtodHRwOi8vd3d3Lm5ldGxvY2suaHUvZG9jcy8wZwYIKwYBBQUHAgIwWwxZTmVtIG1pbsWRc8OtdGV0dCB0YW7DunPDrXR2w6FueS4gU3pvbGfDoWx0YXTDoXNpIHN6YWLDoWx5emF0OiBodHRwOi8vd3d3Lm5ldGxvY2suaHUvZG9jcy8wDQYJKoZIhvcNAQELBQADggEBAIViAXTy8XrP_feYWjU9BWSoaYKZ5MlIM1wi1nKelo9DLpCZUv4VZrXatQIlv4hh3OMxxsKVLJbXCwQTqkuPrPx-VPFf_kmEO0ZUT9aQhJEIm5PDFudYPdu_yROkfpmVseRcG4zyl12yWFIChCHBxQW9ApKhJCitTLo4TLHMtoYnYx8ZQFCCUvxg9pdxQfctm5EQiZfMSFz7uRrJ3uRxSoNpvn39EfjjoYulP-oJlQF_bKwQc7AaflJCs8HjwoqfPPkooyInRNLUPSqajJmXMg5cxEOgHq-7RBXXabPZYU9vhIWlRh9mXc5SGjG8D5QsCS9ywzfcQ831nY4Gbke6M3E`

const signatureB64 = `BIIBAKif89u2VjZs_xk8VEgZsuJjy1VnHTIHNGUfoSE-_ny-TvB2-OTtQUNmVCERGWNhsEUCLn04_fk4xvfM3Cg5h8rZFj2AlBKEZtbb8Ummy4uZLl7erdtN2BKNIQFi-TdFGp4zdpcFlDqx4u5EIEjK6Y4FiYe-XPAXnxcns3KB76mVbSnqg-E5pFEhgCzJz6u4uN_nPAAgukyCxboFmrjj2NonIiWge6qoVFus6nbuiQtQ64NFgeQA7WKlv1jL08P1zUpzEjvttIsTQJUY2Z_zcuqMBd9HMOnTgtQG6R3F6mloAxgOHAGmmL7eKLIvkw38ijeJLGCBJt56l8sVO7JAbBw`

const p7bPEM = `-----BEGIN PKCS7-----
MIIJtAYJKoZIhvcNAQcCoIIJpTCCCaECAQExDzANBglghkgBZQMEAgEFADALBgkq
hkiG9w0BBwGgggcZMIIHFTCCBf2gAwIBAgIOTQkYvAd2pLRlx7SyaLswDQYJKoZI
hvcNAQELBQAwgbUxCzAJBgNVBAYTAkhVMREwDwYDVQQHDAhCdWRhcGVzdDEVMBMG
A1UECgwMTmV0TG9jayBLZnQuMTcwNQYDVQQLDC5UYW7DunPDrXR2w6FueWtpYWTD
s2sgKENlcnRpZmljYXRpb24gU2VydmljZXMpMUMwQQYDVQQDDDpOZXRMb2NrIEV4
cHJlc3N6IEVhdC4gKENsYXNzIEMgTGVnYWwpIFRhbsO6c8OtdHbDoW55a2lhZMOz
MB4XDTEzMTIwNTEzMTAzNFoXDTE1MTIwNTEzMTAzNFowgZIxCzAJBgNVBAYTAkhV
MREwDwYDVQQHDAhCdWRhcGVzdDEdMBsGA1UECgwUTWFneWFyIFRlbGVrb20gTnly
dC4xDjAMBgNVBAsMBVRFU1pUMR0wGwYDVQQDDBRNYWd5YXIgVGVsZWtvbSBOeXJ0
LjEiMCAGCSqGSIb3DQEJARYTcGtpLXRlYW1AdGVsZWtvbS5odTCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAOYLluDTTsyQOqgY6LupOZ1EIdFPV/1rOMWo
vpqxpFXJBLcAnH5Dvt+RyIWzvTOe/ALzdgtPu9C0KjtB7qpzaHbcsTc7PLmqhwOt
HYUacSNBZU2bmK8rBr7r5irIC3+0Q8o+/qMYMqWTY//nXLCo6+aM8+xTIOp5otdZ
jxBK7jpWeFRJerLkjnXFVHSQIPxj7GQn/DsxI97cRPDgyQJcLgEb4OfXu4srIk4y
mrpaIj8YIZ06yBBnhg4a/XC77MVlUIw/rDTMtv2mFotxn6w6Wbt2v9Zg1C1o8ROY
5n8WPQFPjxxqgGzBi7ssxhKIPkp2ab6U5h2DnDs3qxSztgIihckCAwEAAaOCA0Iw
ggM+MAwGA1UdEwEB/wQCMAAwDgYDVR0PAQH/BAQDAgbAMB0GA1UdDgQWBBT1CZ/a
em6MvjoFOXTiI+1ytxxfoTAfBgNVHSMEGDAWgBRxSW6Bh9EIrfEoBKHzv/9Gb3+1
2jAdBgNVHSUEFjAUBggrBgEFBQcDBAYIKwYBBQUHAwIwHgYDVR0RBBcwFYETcGtp
LXRlYW1AdGVsZWtvbS5odTCCAUQGCCsGAQUFBwEBBIIBNjCCATIwLQYIKwYBBQUH
MAGGIWh0dHA6Ly9vY3NwMS5uZXRsb2NrLmh1L2NjbGNhLmNnaTAtBggrBgEFBQcw
AYYhaHR0cDovL29jc3AyLm5ldGxvY2suaHUvY2NsY2EuY2dpMC0GCCsGAQUFBzAB
hiFodHRwOi8vb2NzcDMubmV0bG9jay5odS9jY2xjYS5jZ2kwNQYIKwYBBQUHMAKG
KWh0dHA6Ly9haWExLm5ldGxvY2suaHUvaW5kZXguY2dpP2NhPWNjbGNhMDUGCCsG
AQUFBzAChilodHRwOi8vYWlhMi5uZXRsb2NrLmh1L2luZGV4LmNnaT9jYT1jY2xj
YTA1BggrBgEFBQcwAoYpaHR0cDovL2FpYTMubmV0bG9jay5odS9pbmRleC5jZ2k/
Y2E9Y2NsY2EwgaEGA1UdHwSBmTCBljAwoC6gLIYqaHR0cDovL2NybDEubmV0bG9j
ay5odS9pbmRleC5jZ2k/Y3JsPWNjbGNhMDCgLqAshipodHRwOi8vY3JsMi5uZXRs
b2NrLmh1L2luZGV4LmNnaT9jcmw9Y2NsY2EwMKAuoCyGKmh0dHA6Ly9jcmwzLm5l
dGxvY2suaHUvaW5kZXguY2dpP2NybD1jY2xjYTCBsgYDVR0gBIGqMIGnMIGkBg0r
BgEEAZtjASaJyvBRMIGSMCcGCCsGAQUFBwIBFhtodHRwOi8vd3d3Lm5ldGxvY2su
aHUvZG9jcy8wZwYIKwYBBQUHAgIwWwxZTmVtIG1pbsWRc8OtdGV0dCB0YW7DunPD
rXR2w6FueS4gU3pvbGfDoWx0YXTDoXNpIHN6YWLDoWx5emF0OiBodHRwOi8vd3d3
Lm5ldGxvY2suaHUvZG9jcy8wDQYJKoZIhvcNAQELBQADggEBAIViAXTy8XrP/feY
WjU9BWSoaYKZ5MlIM1wi1nKelo9DLpCZUv4VZrXatQIlv4hh3OMxxsKVLJbXCwQT
qkuPrPx+VPFf/kmEO0ZUT9aQhJEIm5PDFudYPdu/yROkfpmVseRcG4zyl12yWFIC
hCHBxQW9ApKhJCitTLo4TLHMtoYnYx8ZQFCCUvxg9pdxQfctm5EQiZfMSFz7uRrJ
3uRxSoNpvn39EfjjoYulP+oJlQF/bKwQc7AaflJCs8HjwoqfPPkooyInRNLUPSqa
jJmXMg5cxEOgHq+7RBXXabPZYU9vhIWlRh9mXc5SGjG8D5QsCS9ywzfcQ831nY4G
bke6M3ExggJfMIICWwIBATCByDCBtTELMAkGA1UEBhMCSFUxETAPBgNVBAcMCEJ1
ZGFwZXN0MRUwEwYDVQQKDAxOZXRMb2NrIEtmdC4xNzA1BgNVBAsMLlRhbsO6c8Ot
dHbDoW55a2lhZMOzayAoQ2VydGlmaWNhdGlvbiBTZXJ2aWNlcykxQzBBBgNVBAMM
Ok5ldExvY2sgRXhwcmVzc3ogRWF0LiAoQ2xhc3MgQyBMZWdhbCkgVGFuw7pzw610
dsOhbnlraWFkw7MCDk0JGLwHdqS0Zce0smi7MA0GCWCGSAFlAwQCAQUAoGkwGAYJ
KoZIhvcNAQkDMQsGCSqGSIb3DQEHATAcBgkqhkiG9w0BCQUxDxcNMTUwMjA2MTAy
MjM3WjAvBgkqhkiG9w0BCQQxIgQghqzGKYxwixbMJNPwlgjd/fwkwaV8PM+TjU1t
8JE//JwwDQYJKoZIhvcNAQEBBQAEggEAqJ/z27ZWNmz/GTxUSBmy4mPLVWcdMgc0
ZR+hIT7+fL5O8Hb45O1BQ2ZUIREZY2GwRQIufTj9+TjG98zcKDmHytkWPYCUEoRm
1tvxSabLi5kuXt6t203YEo0hAWL5N0UanjN2lwWUOrHi7kQgSMrpjgWJh75c8Bef
FyezcoHvqZVtKeqD4TmkUSGALMnPq7i43+c8ACC6TILFugWauOPY2iciJaB7qqhU
W6zqdu6JC1Drg0WB5ADtYqW/WMvTw/XNSnMSO+20ixNAlRjZn/Ny6owF30cw6dOC
1AbpHcXqaWgDGA4cAaaYvt4osi+TDfyKN4ksYIEm3nqXyxU7skBsHA==
-----END PKCS7-----`
