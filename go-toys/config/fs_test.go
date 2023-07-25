package config

type fsProvider struct {
	basePath string
}

var _ Provider = &fsProvider{}

func (p *fsProvider) Actual() (State, error) {
	panic("implement me")
}

func (p *fsProvider) LookupVersion(checksum string) (State, error) {
	//TODO implement me
	panic("implement me")
}
