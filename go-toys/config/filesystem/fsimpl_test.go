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

func TestInheritance(t *testing.T) {
	type A struct {
		name string
	}
	type B struct {
		A
		size int
	}

	b := &B{
		A: A{
			name: "foo",
		},
		size: 1,
	}

	t.Logf("%s, %d", b.name, b.size)

	a := b.A

	t.Logf("%s", a.name)

}
