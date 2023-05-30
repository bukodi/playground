package config

import (
	"github.com/knadh/koanf/v2"
	"testing"
)

type State interface {
	Checksum() string
	Validate() error
	GetBytes(path string) ([]byte, error)
	GetBytesWithDefault(path string, defValue []byte) []byte
	GetPaths(basePath string) ([]string, error)
}

func TestKoanf(t *testing.T) {
	k := koanf.NewWithConf(koanf.Conf{
		Delim: ".",
	})
	_ = k
}
