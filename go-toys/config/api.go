package config

type Provider interface {
	Actual() (State, error)
	LookupVersion(checksum string) (State, error)
}

type State interface {
	Checksum() string
	Validate() error
	GetBytes(path string) ([]byte, error)
	GetBytesWithDefault(path string, defValue []byte) []byte
	GetPaths(basePath string) ([]string, error)
}
