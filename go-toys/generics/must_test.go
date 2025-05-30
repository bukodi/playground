package generics

import (
	"encoding/base64"
	"testing"
)

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func TestMust(t *testing.T) {

	data, err := base64.StdEncoding.DecodeString("aGVsbG8gd29ybGQ=")
	if err != nil {
		t.Fatalf("%+v", err)
	} else {
		t.Logf("data: %s", data)
	}

	data2 := Must(base64.StdEncoding.DecodeString("aGVsbG8gd29ybGQ="))
	t.Logf("data2: %s", data2)

}
