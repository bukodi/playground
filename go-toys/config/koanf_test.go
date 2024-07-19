package config

import (
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"testing"
)

type myProvider struct {
	parent koanf.Provider
}

func (myProv myProvider) ReadBytes() ([]byte, error) {
	return myProv.parent.ReadBytes()
}

func (myProv myProvider) Read() (map[string]interface{}, error) {
	return myProv.parent.Read()
}

var _ koanf.Provider = (*myProvider)(nil)

type myParser struct {
	parent koanf.Parser
}

func (myParser myParser) Unmarshal(bytes []byte) (map[string]interface{}, error) {
	return myParser.parent.Unmarshal(bytes)
}

func (myParser myParser) Marshal(m2 map[string]interface{}) ([]byte, error) {
	return myParser.parent.Marshal(m2)
}

var _ koanf.Parser = (*myParser)(nil)

func TestKoanf(t *testing.T) {
	// Define the config struct
	type Config struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
	}

	// Define the config file
	configFile := []byte(`{
		"host": "localhost",
		"port": 8080
	}`)

	rawbytes.Provider(configFile)

	k := koanf.New(".")
	k.Load(rawbytes.Provider(configFile), json.Parser())
}
