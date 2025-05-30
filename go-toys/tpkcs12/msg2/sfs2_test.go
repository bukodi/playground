package msg2

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"testing"
)

func TestRSASign(t *testing.T) {
	sk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	msg := []byte("hello, world")
	digest := sha256.Sum256(msg)
	t.Logf("digest: %x", digest)
	sign, err := rsa.SignPKCS1v15(rand.Reader, sk, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("sign: %x", sign)
	t.Logf("len(sign): %d", len(sign))

	err = rsa.VerifyPKCS1v15(&sk.PublicKey, crypto.SHA256, digest[:], sign)
	if err != nil {
		t.Fatalf("Verification failed: %+v", err)
	} else {
		t.Logf("signature verified")
	}

}

func TestYakud(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("MS43LjAgW0J1aWxkOiA0MzRdADEuMABwN3BlbWRldAAAAA==")
	if err != nil {
		t.Fatalf("%+v", err)
	} else {
		t.Logf("data: %x", data)
	}
	digest := sha256.Sum256(data)
	t.Logf("digest: %x", digest)
}

func TestInput2(t *testing.T) {
	inputData, err := base64.StdEncoding.DecodeString(inputB64)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	sha256Bytes := sha256.Sum256(inputData)
	t.Logf("input as hex: %s", hex.EncodeToString(inputData))
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
	} else if bytes.Equal(p7Bytes[60:60+2118], embeddedCertBytes) {
		t.Logf("embeddedCertBytes match")
	} else {
		t.Logf("embeddedCertBytes do not match")
	}
	embeddedCert, err := x509.ParseCertificate(embeddedCertBytes)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("embeddedCert.Subject.CommonName: %s", embeddedCert.Subject.CommonName)

	signedAttrsBytes, err := base64.RawURLEncoding.DecodeString(signedAttrsB64)
	if err != nil {
		t.Errorf("%+v", err)
	} else if bytes.Equal(p7Bytes[2320:2320+4+528], signedAttrsBytes) {
		t.Logf("signedAttrsBytes match")
	} else {
		t.Fatalf("signedAttrsBytes do not match")
	}
	//signedAttrsBytes = signedAttrsBytes[4:]
	signedAttrsBytes[0] = 0x31
	signedAttrsHash := sha256.Sum256(signedAttrsBytes)
	//signedAttrsBytes = signedAttrsBytes[2:]
	/*msgDigestAttrBytes, err := base64.RawURLEncoding.DecodeString(msgDigestAttrB64)
	if err != nil {
		t.Errorf("%+v", err)
	} else if bytes.Equal(p7Bytes[2164:2164+49], msgDigestAttrBytes) {
		t.Logf("msgDigestAttrBytes match")
	} else {
		t.Fatalf("msgDigestAttrBytes do not match")
	}
	msgDigestAttrBytes = msgDigestAttrBytes[4:]
	signedAttrsHash := sha256.Sum256(msgDigestAttrBytes)*/

	signatureBytes, err := base64.RawURLEncoding.DecodeString(signatureB64)
	signatureBytes = signatureBytes[4:]
	if err != nil {
		t.Errorf("%+v", err)
	} else if bytes.Equal(p7Bytes[2867+4:2867+4+256], signatureBytes) {
		t.Logf("signatureBytes match")
	} else {
		t.Fatalf("signatureBytes do not match")
	}

	_ = signedAttrsHash
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

const inputB64 = "MS4xMC42ODEgW0J1aWxkOiAwODkwXQAxLjEAcDdwZW1kZXQAAHRrb20AZ3NfcXRrb20yMDE3AAA="

const signedAttrsB64 = `oIICEDAWBgcqhkiG9w0CMQsGCWCGSAFlAwQCATAWBgcqhkiG9w0DMQsGCSqGSIb3DQEBATAYBgkqhkiG9w0BCQMxCwYJKoZIhvcNAQcBMBwGCSqGSIb3DQEJBTEPFw0yNTAzMjUxMzA2NDNaMC0GCSqGSIb3DQEJNDEgMB4wDQYJYIZIAWUDBAIBBQChDQYJKoZIhvcNAQEBBQAwLwYJKoZIhvcNAQcFMSIEIMNyRh9nOgOSPJX-z1UVqBJlkp2H3kTe0zw7uUJuv3N6MC8GCSqGSIb3DQEJBDEiBCDDckYfZzoDkjyV_s9VFagSZZKdh95E3tM8O7lCbr9zejCCARMGCSqGSIb3DQEHBjGCAQQEggEAbIml5ySO37c0_EYZx-h4940uwNZC2w_lZ_8n-hY3fUjKCKTtwFpQIzPmg6k4nENaSD-7f64fV215Vx6L8gxD-Cn4k6yovhEv1LeAqXq5W_5HbZXa4TZu0vk7onpoTeW1rEb631qJ4J1f2h9fnfMtlP-V3sjGkVv9AYR_5bhcmg_BPRZGMToFppuLHIZHQlPDM6RBgenJh9I0DPT9P3y9yyfTGTwVX9Cn5YDB6uxdKCGyByrvAcAhrLDm7wmqqU3uyyNr4C0nAyRUYoodGHgNU3xDjcYT8oBj3ePQHVQx82MTQ8vniSQdNVApiQ0dckl2UmWK4DL5ALj8qv6GLVxsIw`

const msgDigestAttrB64 = `MC8GCSqGSIb3DQEJBDEiBCCGrMYpjHCLFswk0_CWCN39_CTBpXw8z5ONTW3wkT_8nA`

const embeddedCertB64 = `MIIIQjCCByqgAwIBAgIOV-vxwwUCxEannbMshoswDQYJKoZIhvcNAQELBQAwYDELMAkGA1UEBhMCSFUxETAPBgNVBAcMCEJ1ZGFwZXN0MRUwEwYDVQQKDAxORVRMT0NLIEx0ZC4xJzAlBgNVBAMMHk5FVExPQ0sgVHJ1c3QgUXVhbGlmaWVkIFNDRCBDQTAeFw0yMzA0MjcxNDU5NDNaFw0yNTA0MjYxNDU5NDNaMIHPMQswCQYDVQQGEwJIVTERMA8GA1UEBwwIQnVkYXBlc3QxHTAbBgNVBAoMFE1hZ3lhciBUZWxla29tIE55cnQuMR0wGwYDVQQDDBRNYWd5YXIgVGVsZWtvbSBOeXJ0LjEtMCsGA1UEBRMkMS4zLjYuMS40LjEuMzU1NS41LjIuNTI3MTUyODI4MDA4MDk4MSIwIAYJKoZIhvcNAQkBFhNwa2ktdGVhbUB0ZWxla29tLmh1MRwwGgYDVQRhDBNWQVRIVS0xMDc3MzM4MS0yLTQ0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0gY0zkSt7_RHeVFudS9LaxATmBlNbA-4cZKGHDCL0H3aO3_Vd3mJcKUhs4n6nCM6pIS9seDzqFCjNh99SRXAqkuQU6OdBWxYq8Y8oElLiQ3x52fiYTv-dX41V1ccTD0y8EZfc41dpJjZFS-xJNGXxeCsaRSxnXqfkRDZZw5SycERgaDvZaGVGvgnlTuXviE7gUoAIXl_4MXVbormzS4B3cEKQJbYhJesNEcAoDqnYhTC5QxbOnU0118jaUmYHSCoPWxrkr9C1DmqlXiPrxXOKvzj5VgLqLt9XQIiH6aSzzItoWT5yUnGc2thscJQkGaICKkmdaasD_ezyNxhBxAZYQIDAQABo4IEiDCCBIQwDAYDVR0TAQH_BAIwADAOBgNVHQ8BAf8EBAMCBkAwgacGCCsGAQUFBwEDBIGaMIGXMAgGBgQAjkYBATAVBgYEAI5GAQIwCxMDSFVGAgEFAgEGMAsGBgQAjkYBAwIBCjBSBgYEAI5GAQUwSDAiFhxodHRwczovL3d3dy5uZXRsb2NrLmh1L2RvY3MvEwJlbjAiFhxodHRwczovL3d3dy5uZXRsb2NrLmh1L2RvY3MvEwJodTATBgYEAI5GAQYwCQYHBACORgEGAjCCAQ0GA1UdIASCAQQwggEAMBAGDisGAQQBgvcQAwICAgICMIHrBgcEAIvsQAEBMIHfMCcGCCsGAQUFBwIBFhtodHRwOi8vd3d3Lm5ldGxvY2suaHUvZG9jcy8wgbMGCCsGAQUFBwICMIGmDIGjUXVhbGlmaWVkIGNlcnRpZmljYXRlLCBpc3N1ZWQgdG8gbGVnYWwgcGVyc29uLCBmb3IgY3JlYXRpbmcgYWR2YW5jZWQgc2VhbC4gSm9naSBzemVtw6lseW5layBraWFkb3R0IG1pbsWRc8OtdGV0dCB0YW7DunPDrXR2w6FueSBmb2tvem90dCBiw6lseWVnIGzDqXRyZWhvesOhc8OhaG96LjAdBgNVHQ4EFgQUYCIx4e7RiPQsnoWvHOl_lwJdTJEwHwYDVR0jBBgwFoAUVH7spB4YW7vXs3mD9PvXZ7splfwwggFcBggrBgEFBQcBAQSCAU4wggFKMDEGCCsGAQUFBzABhiVodHRwOi8vb2NzcDEubmV0bG9jay5odS9xdHJ1c3RzY2QuY2dpMDEGCCsGAQUFBzABhiVodHRwOi8vb2NzcDIubmV0bG9jay5odS9xdHJ1c3RzY2QuY2dpMDEGCCsGAQUFBzABhiVodHRwOi8vb2NzcDMubmV0bG9jay5odS9xdHJ1c3RzY2QuY2dpMDkGCCsGAQUFBzAChi1odHRwOi8vYWlhMS5uZXRsb2NrLmh1L2luZGV4LmNnaT9jYT1xdHJ1c3RzY2QwOQYIKwYBBQUHMAKGLWh0dHA6Ly9haWEyLm5ldGxvY2suaHUvaW5kZXguY2dpP2NhPXF0cnVzdHNjZDA5BggrBgEFBQcwAoYtaHR0cDovL2FpYTMubmV0bG9jay5odS9pbmRleC5jZ2k_Y2E9cXRydXN0c2NkMIGtBgNVHR8EgaUwgaIwNKAyoDCGLmh0dHA6Ly9jcmwxLm5ldGxvY2suaHUvaW5kZXguY2dpP2NybD1xdHJ1c3RzY2QwNKAyoDCGLmh0dHA6Ly9jcmwyLm5ldGxvY2suaHUvaW5kZXguY2dpP2NybD1xdHJ1c3RzY2QwNKAyoDCGLmh0dHA6Ly9jcmwzLm5ldGxvY2suaHUvaW5kZXguY2dpP2NybD1xdHJ1c3RzY2QwOAYDVR0RBDEwL4ETcGtpLXRlYW1AdGVsZWtvbS5odaAYBggrBgEFBQcIA6AMMAoGCCsGAQQBm2MFMB8GA1UdJQQYMBYGCisGAQQBgjcKAwwGCCsGAQUFBwMEMA0GCSqGSIb3DQEBCwUAA4IBAQDH-PDJtt8643BwptFUYE2ciXpfi7Mewr5YjEUDz9jeTRRPOKngZfo5f86-vc2CakAP00BS3B5rW6JBMzdVP4afwCcuf2piw4DogYenW3g_iyYIPLZB4lzHi2s-m8mkjUEXO2ZVQcDkMK5f52NpuTiEnxqMcZjyQYif1VTedah58qP2CLKDnpod-_-KsrvvWyYZ94ki_Y3Tmx99kpmzKU_niTjz9DC3wkkUDAzPHnVnHngVPCGPZ1aIxIp5BPaxJLoKU58miVlUf0fUjqaobunjROBuMEf7qdhYT_PtgXLY-Lr8rsIQz9zfV1XtgNsBgWRUsovqg6yKlpuMg0mvJpXx`

const signatureB64 = `BIIBACt3XXJBIPb7h0LGFKds5Gk1bGPjb5prO9Ci_eAi-sEV7AyOf6nyO0K_Jax-ek6fakbOvvFYOeFdesgcsB9a8qrgB6ErgT5HjDaKxaMqNoiH0w1UjN7gMEH5MqD84a6mwQzwK6ESG3TSQ7WbV2184-9GL0WaKNeEHCk8RxLOns_4SuLHw0ahUpz-8DIHa705g98KXDtaaaY_vnUvu5vB7_45oUWUBFVxqybs3wCfLHrJAqFqfCMaxhmHqUm5SAsOWDSeKMzZQ4BOQPFGE7UK5Y_TuYhZ2bE-CeL6Wc1QRdsrxFsbVl8VR8tT0q9WTDFmgbllwwuunc3YVh5RbuoPXV0`

const p7bPEM = `-----BEGIN PKCS7-----
MIIMMwYJKoZIhvcNAQcCoIIMJDCCDCACAQExDzANBglghkgBZQMEAgEFADALBgkq
hkiG9w0BBwGggghGMIIIQjCCByqgAwIBAgIOV+vxwwUCxEannbMshoswDQYJKoZI
hvcNAQELBQAwYDELMAkGA1UEBhMCSFUxETAPBgNVBAcMCEJ1ZGFwZXN0MRUwEwYD
VQQKDAxORVRMT0NLIEx0ZC4xJzAlBgNVBAMMHk5FVExPQ0sgVHJ1c3QgUXVhbGlm
aWVkIFNDRCBDQTAeFw0yMzA0MjcxNDU5NDNaFw0yNTA0MjYxNDU5NDNaMIHPMQsw
CQYDVQQGEwJIVTERMA8GA1UEBwwIQnVkYXBlc3QxHTAbBgNVBAoMFE1hZ3lhciBU
ZWxla29tIE55cnQuMR0wGwYDVQQDDBRNYWd5YXIgVGVsZWtvbSBOeXJ0LjEtMCsG
A1UEBRMkMS4zLjYuMS40LjEuMzU1NS41LjIuNTI3MTUyODI4MDA4MDk4MSIwIAYJ
KoZIhvcNAQkBFhNwa2ktdGVhbUB0ZWxla29tLmh1MRwwGgYDVQRhDBNWQVRIVS0x
MDc3MzM4MS0yLTQ0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0gY0
zkSt7/RHeVFudS9LaxATmBlNbA+4cZKGHDCL0H3aO3/Vd3mJcKUhs4n6nCM6pIS9
seDzqFCjNh99SRXAqkuQU6OdBWxYq8Y8oElLiQ3x52fiYTv+dX41V1ccTD0y8EZf
c41dpJjZFS+xJNGXxeCsaRSxnXqfkRDZZw5SycERgaDvZaGVGvgnlTuXviE7gUoA
IXl/4MXVbormzS4B3cEKQJbYhJesNEcAoDqnYhTC5QxbOnU0118jaUmYHSCoPWxr
kr9C1DmqlXiPrxXOKvzj5VgLqLt9XQIiH6aSzzItoWT5yUnGc2thscJQkGaICKkm
daasD/ezyNxhBxAZYQIDAQABo4IEiDCCBIQwDAYDVR0TAQH/BAIwADAOBgNVHQ8B
Af8EBAMCBkAwgacGCCsGAQUFBwEDBIGaMIGXMAgGBgQAjkYBATAVBgYEAI5GAQIw
CxMDSFVGAgEFAgEGMAsGBgQAjkYBAwIBCjBSBgYEAI5GAQUwSDAiFhxodHRwczov
L3d3dy5uZXRsb2NrLmh1L2RvY3MvEwJlbjAiFhxodHRwczovL3d3dy5uZXRsb2Nr
Lmh1L2RvY3MvEwJodTATBgYEAI5GAQYwCQYHBACORgEGAjCCAQ0GA1UdIASCAQQw
ggEAMBAGDisGAQQBgvcQAwICAgICMIHrBgcEAIvsQAEBMIHfMCcGCCsGAQUFBwIB
FhtodHRwOi8vd3d3Lm5ldGxvY2suaHUvZG9jcy8wgbMGCCsGAQUFBwICMIGmDIGj
UXVhbGlmaWVkIGNlcnRpZmljYXRlLCBpc3N1ZWQgdG8gbGVnYWwgcGVyc29uLCBm
b3IgY3JlYXRpbmcgYWR2YW5jZWQgc2VhbC4gSm9naSBzemVtw6lseW5layBraWFk
b3R0IG1pbsWRc8OtdGV0dCB0YW7DunPDrXR2w6FueSBmb2tvem90dCBiw6lseWVn
IGzDqXRyZWhvesOhc8OhaG96LjAdBgNVHQ4EFgQUYCIx4e7RiPQsnoWvHOl/lwJd
TJEwHwYDVR0jBBgwFoAUVH7spB4YW7vXs3mD9PvXZ7splfwwggFcBggrBgEFBQcB
AQSCAU4wggFKMDEGCCsGAQUFBzABhiVodHRwOi8vb2NzcDEubmV0bG9jay5odS9x
dHJ1c3RzY2QuY2dpMDEGCCsGAQUFBzABhiVodHRwOi8vb2NzcDIubmV0bG9jay5o
dS9xdHJ1c3RzY2QuY2dpMDEGCCsGAQUFBzABhiVodHRwOi8vb2NzcDMubmV0bG9j
ay5odS9xdHJ1c3RzY2QuY2dpMDkGCCsGAQUFBzAChi1odHRwOi8vYWlhMS5uZXRs
b2NrLmh1L2luZGV4LmNnaT9jYT1xdHJ1c3RzY2QwOQYIKwYBBQUHMAKGLWh0dHA6
Ly9haWEyLm5ldGxvY2suaHUvaW5kZXguY2dpP2NhPXF0cnVzdHNjZDA5BggrBgEF
BQcwAoYtaHR0cDovL2FpYTMubmV0bG9jay5odS9pbmRleC5jZ2k/Y2E9cXRydXN0
c2NkMIGtBgNVHR8EgaUwgaIwNKAyoDCGLmh0dHA6Ly9jcmwxLm5ldGxvY2suaHUv
aW5kZXguY2dpP2NybD1xdHJ1c3RzY2QwNKAyoDCGLmh0dHA6Ly9jcmwyLm5ldGxv
Y2suaHUvaW5kZXguY2dpP2NybD1xdHJ1c3RzY2QwNKAyoDCGLmh0dHA6Ly9jcmwz
Lm5ldGxvY2suaHUvaW5kZXguY2dpP2NybD1xdHJ1c3RzY2QwOAYDVR0RBDEwL4ET
cGtpLXRlYW1AdGVsZWtvbS5odaAYBggrBgEFBQcIA6AMMAoGCCsGAQQBm2MFMB8G
A1UdJQQYMBYGCisGAQQBgjcKAwwGCCsGAQUFBwMEMA0GCSqGSIb3DQEBCwUAA4IB
AQDH+PDJtt8643BwptFUYE2ciXpfi7Mewr5YjEUDz9jeTRRPOKngZfo5f86+vc2C
akAP00BS3B5rW6JBMzdVP4afwCcuf2piw4DogYenW3g/iyYIPLZB4lzHi2s+m8mk
jUEXO2ZVQcDkMK5f52NpuTiEnxqMcZjyQYif1VTedah58qP2CLKDnpod+/+Ksrvv
WyYZ94ki/Y3Tmx99kpmzKU/niTjz9DC3wkkUDAzPHnVnHngVPCGPZ1aIxIp5BPax
JLoKU58miVlUf0fUjqaobunjROBuMEf7qdhYT/PtgXLY+Lr8rsIQz9zfV1XtgNsB
gWRUsovqg6yKlpuMg0mvJpXxMYIDsTCCA60CAQEwcjBgMQswCQYDVQQGEwJIVTER
MA8GA1UEBwwIQnVkYXBlc3QxFTATBgNVBAoMDE5FVExPQ0sgTHRkLjEnMCUGA1UE
AwweTkVUTE9DSyBUcnVzdCBRdWFsaWZpZWQgU0NEIENBAg5X6/HDBQLERqedsyyG
izANBglghkgBZQMEAgEFAKCCAhAwFgYHKoZIhvcNAjELBglghkgBZQMEAgEwFgYH
KoZIhvcNAzELBgkqhkiG9w0BAQEwGAYJKoZIhvcNAQkDMQsGCSqGSIb3DQEHATAc
BgkqhkiG9w0BCQUxDxcNMjUwMzI1MTMwNjQzWjAtBgkqhkiG9w0BCTQxIDAeMA0G
CWCGSAFlAwQCAQUAoQ0GCSqGSIb3DQEBAQUAMC8GCSqGSIb3DQEHBTEiBCDDckYf
ZzoDkjyV/s9VFagSZZKdh95E3tM8O7lCbr9zejAvBgkqhkiG9w0BCQQxIgQgw3JG
H2c6A5I8lf7PVRWoEmWSnYfeRN7TPDu5Qm6/c3owggETBgkqhkiG9w0BBwYxggEE
BIIBAGyJpeckjt+3NPxGGcfoePeNLsDWQtsP5Wf/J/oWN31Iygik7cBaUCMz5oOp
OJxDWkg/u3+uH1dteVcei/IMQ/gp+JOsqL4RL9S3gKl6uVv+R22V2uE2btL5O6J6
aE3ltaxG+t9aieCdX9ofX53zLZT/ld7IxpFb/QGEf+W4XJoPwT0WRjE6BaabixyG
R0JTwzOkQYHpyYfSNAz0/T98vcsn0xk8FV/Qp+WAwersXSghsgcq7wHAIayw5u8J
qqlN7ssja+AtJwMkVGKKHRh4DVN8Q43GE/KAY93j0B1UMfNjE0PL54kkHTVQKYkN
HXJJdlJliuAy+QC4/Kr+hi1cbCMwDQYJKoZIhvcNAQEBBQAEggEAK3ddckEg9vuH
QsYUp2zkaTVsY+Nvmms70KL94CL6wRXsDI5/qfI7Qr8lrH56Tp9qRs6+8Vg54V16
yBywH1ryquAHoSuBPkeMNorFoyo2iIfTDVSM3uAwQfkyoPzhrqbBDPAroRIbdNJD
tZtXbXzj70YvRZoo14QcKTxHEs6ez/hK4sfDRqFSnP7wMgdrvTmD3wpcO1pppj++
dS+7m8Hv/jmhRZQEVXGrJuzfAJ8seskCoWp8IxrGGYepSblICw5YNJ4ozNlDgE5A
8UYTtQrlj9O5iFnZsT4J4vpZzVBF2yvEWxtWXxVHy1PSr1ZMMWaBuWXDC66dzdhW
HlFu6g9dXQ==
-----END PKCS7-----`
