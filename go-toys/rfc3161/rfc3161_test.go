package rfc3161

import (
	"crypto"
	"crypto/sha1"
	"os"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	req, err := ReadTSQ("./testdata/sha1.tsq")
	if err != nil {
		t.Error(err)
	}
	err = req.Verify()
	if err != nil {
		t.Error(err)
	}
	_, err = ReadTSR("./testdata/sha1.response.tsr")
	if err != nil {
		t.Error(err)
	}

	req, err = ReadTSQ("./testdata/sha1_nonce.tsq")
	if err != nil {
		t.Error(err)
	}
	err = req.Verify()
	if err != nil {
		t.Error(err)
	}
	_, err = ReadTSR("./testdata/sha1_nonce.response.tsr")
	if err != nil {
		t.Error(err)
	}
}

// Contruct the tsr manually
func TestReqBuildManually(t *testing.T) {
	mes, err := os.ReadFile("./testdata/message.txt")
	if err != nil {
		t.Error(err)
	}
	digest := sha1.Sum(mes)

	tsr2, err := NewTimeStampReq(crypto.SHA1, digest[:])
	if err != nil {
		t.Error(err)
	}
	err = tsr2.GenerateNonce()
	if err != nil {
		t.Error(err)
	}
	err = tsr2.Verify()
	if err != nil {
		t.Error(err)
	}
}
