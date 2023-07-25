package config

// Empty directories and files will be skipped
// hidden files and directories can be saved, but its excluded from checksum calculation

type Provider interface {
	Actual() (State, error)
	LookupVersion(checksum string) (State, error)
}

type State interface {
	Checksum() string
	Size() uint
	Get(path string) []byte
	List(basePath string) map[string][]byte
	Add(files map[string][]byte) State
}
