package config

import (
	"os"
	"path/filepath"
)

type fsProvider struct {
	basePath string
}

var _ Provider = &fsProvider{}

func (p *fsProvider) Actual() (State, error) {
	s := fsState{
		basePath: p.basePath,
	}
	return &s, nil
}

func (p *fsProvider) LookupVersion(checksum string) (State, error) {
	//TODO implement me
	panic("implement me")
}

type fsState struct {
	basePath        string
	cached_checksum string
}

func (f *fsState) Checksum() string {
	if f.cached_checksum != "" {
		return f.cached_checksum
	}

	//TODO implement me
	panic("implement me")
}

func (f *fsState) Validate() error {
	//TODO implement me
	panic("implement me")
}

func (f *fsState) GetBytes(path string) ([]byte, error) {
	fullPath := filepath.Join(f.basePath, path)
	if fileContent, err := os.ReadFile(fullPath); err != nil {
		return nil, err
	} else {
		return fileContent, nil
	}
}

func (f *fsState) GetBytesWithDefault(path string, defValue []byte) []byte {
	if content, err := f.GetBytes(path); err == nil {
		return content
	} else {
		return defValue
	}
}

func (f *fsState) GetPaths(basePath string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

var _ State = &fsState{}
