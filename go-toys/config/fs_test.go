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
	panic("implement me")
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

func (f *fsState) GetPaths(path string) (keyNames []string, subDirNames []string, err error) {
	fullPath := filepath.Join(f.basePath, path)
	dirEntries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, nil, err
	}
	for _, de := range dirEntries {
		if de.IsDir() {
			subDirNames = append(subDirNames, de.Name())
		} else {
			keyNames = append(keyNames, de.Name())
		}
	}
	return keyNames, subDirNames, nil
}
