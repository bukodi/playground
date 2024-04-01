package errlog

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"testing"
)

// SensitiveString is a string that should not be logged
type SensitiveString string

func (ss SensitiveString) String() string {
	return "***Sensitive***"
}

func (ss SensitiveString) GoString() string {
	return "***Sensitive***"
}

func TestSensitiveString(t *testing.T) {
	var ss SensitiveString = "s3cr3t"
	result := fmt.Sprintf("%s, %v, %+v, %#v", ss, ss, ss, ss)
	t.Logf("Result: %s", result)
	t.Logf("Real content: %s", string(ss))
}

func TestSensitiveStringInStruct(t *testing.T) {
	type User struct {
		Username string
		Password SensitiveString
	}

	u := User{
		Username: "Alice",
		Password: "s3cr3t",
	}

	slog.Info("Ok", "user", u)

	result := fmt.Sprintf("%s, %v, %+v, %#v", u, u, u, u)
	t.Logf("Result: %s", result)
	t.Logf("Real content: %s", string(u.Password))
}

type SensitiveRSAKey rsa.PrivateKey

func (sk SensitiveRSAKey) String() string {
	return fmt.Sprintf("Private hidden, public: %v", sk.PublicKey)
}

func (sk SensitiveRSAKey) GoString() string {
	return fmt.Sprintf("%#v", struct {
		Private string
		Public  rsa.PublicKey
	}{
		Private: "***Sensitive***",
		Public:  sk.PublicKey,
	})
}

func TestSensitiveRSAKey(t *testing.T) {
	rk, _ := rsa.GenerateKey(rand.Reader, 2048)
	pk := (*SensitiveRSAKey)(rk)
	t.Logf("Private key summary: %v", pk)
	t.Logf("Private key details: %#v", pk)
}
