package generics

type Key interface {
	Id() [32]byte
}

type KeyStore[T Key] interface {
	Keys() []T
}

type P11Key struct {
}

func (p *P11Key) Id() [32]byte {
	panic("implement me")
}

var _ Key = &P11Key{}

type P11KeyStore struct {
}

func (p11Ks *P11KeyStore) Keys() []*P11Key {
	return nil
}

var _ KeyStore[*P11Key] = &P11KeyStore{}

var Registry = RegistryType{}

type RegistryType struct {
	stores []KeyStore[Key]
}
