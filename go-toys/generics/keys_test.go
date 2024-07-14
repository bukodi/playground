package generics

import "testing"

func TestKs(t *testing.T) {
	p11Ks := &P11KeyStore{}
	for _, key := range p11Ks.Keys() {
		_ = key
	}

}
