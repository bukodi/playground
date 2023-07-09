package config

import (
	"github.com/knadh/koanf/v2"
	"testing"
)

func TestKoanf(t *testing.T) {
	k := koanf.NewWithConf(koanf.Conf{
		Delim: ".",
	})
	_ = k
}
