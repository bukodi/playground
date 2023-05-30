package filesystem

import (
	"crypto/sha256"
	"testing"
)

func TestHashSum(t *testing.T) {

	{
		h := sha256.New()
		h.Write([]byte{1, 2, 3})
		h.Write([]byte{4, 5, 6})
		t.Logf("%x", h.Sum(nil))
	}
	{
		h := sha256.New()
		h.Write([]byte{1, 2, 3, 4, 5, 6})
		t.Logf("%x", h.Sum(nil))
	}
}
