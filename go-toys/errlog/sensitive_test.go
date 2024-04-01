package errlog

import (
	"fmt"
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
	var ss SensitiveString = "password"

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
		Password: "password",
	}

	result := fmt.Sprintf("%s, %v, %+v, %#v", u, u, u, u)
	t.Logf("Result: %s", result)
	t.Logf("Real content: %s", string(u.Password))
}
