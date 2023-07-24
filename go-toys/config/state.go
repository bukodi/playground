package config

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type cfgDir struct {
	subDirs map[string]*cfgDir
	files   map[string][]byte
}

func (d *cfgDir) clone() *cfgDir {
	newD := &cfgDir{}
	if len(d.files) > 0 {
		newD.files = make(map[string][]byte)
		for n, f := range d.files {
			newD.files[n] = f
		}
	}
	if len(d.subDirs) > 0 {
		newD.subDirs = make(map[string]*cfgDir)
		for n, subDir := range d.subDirs {
			newD.subDirs[n] = subDir.clone()
		}
	}
	return newD
}

func (d *cfgDir) find(pathPaths []string) (*cfgDir, []byte) {
	if len(pathPaths) == 0 {
		return d, nil
	}
	if len(pathPaths) == 1 {
		return nil, d.files[pathPaths[0]]
	}
	subDir := d.subDirs[pathPaths[0]]
	if subDir == nil {
		return nil, nil
	}
	return subDir.find(pathPaths[1:])
}

func (d *cfgDir) list(pathPrefix string, ret map[string][]byte) {
	for name, content := range d.files {
		ret[pathPrefix+"/"+name] = content
	}
	for name, subDir := range d.subDirs {
		subDir.list(pathPrefix+"/"+name, ret)
	}
	return
}

func (d *cfgDir) replaceFile(pathParts []string, newContent []byte, replace bool) (content []byte) {
	if len(pathParts) == 1 {
		content = d.files[pathParts[0]]
		if replace {
			if len(newContent) == 0 {
				delete(d.files, pathParts[0])
			} else {
				if d.files == nil {
					d.files = make(map[string][]byte)
				}
				d.files[pathParts[0]] = newContent
			}
		}
		return content
	}

	subDir := d.subDirs[pathParts[0]]
	if subDir == nil && replace {
		subDir = &cfgDir{}
		if d.subDirs == nil {
			d.subDirs = make(map[string]*cfgDir)
		}
		d.subDirs[pathParts[0]] = subDir
	}
	if subDir != nil || replace {
		content = subDir.replaceFile(pathParts[1:], newContent, replace)
	}
	if replace {
		if len(subDir.files) == 0 && len(subDir.subDirs) == 0 {
			delete(d.subDirs, pathParts[0])
		}
	}
	return content
}

type cfgState struct {
	checksum []byte
	size     uint
	rootDir  *cfgDir
}

func (s *cfgState) clone() *cfgState {
	newS := &cfgState{
		size:    s.size,
		rootDir: nil,
	}
	newS.checksum = make([]byte, len(s.checksum))
	copy(newS.checksum, s.checksum)
	newS.rootDir = s.rootDir.clone()
	return newS
}

func (s *cfgState) Size() uint {
	return s.size
}

func (s *cfgState) Get(path string) []byte {
	pathParts := deleteEmpty(strings.Split(path, "/"))
	return s.rootDir.replaceFile(pathParts, nil, false)
}

func (s *cfgState) List(basePath string) map[string][]byte {
	m := make(map[string][]byte)
	pathParts := deleteEmpty(strings.Split(basePath, "/"))
	baseDir, content := s.rootDir.find(pathParts)
	if len(content) > 0 {
		m[basePath] = content
		return m
	}
	if baseDir != nil {
		baseDir.list("", m)
	}
	return m
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func checksum(pathString string, isHidden bool, content []byte) []byte {
	if isHidden || len(content) == 0 {
		return make([]byte, 32)
	}
	hash := sha256.New()
	hash.Write([]byte(pathString))
	hash.Write(content)
	return hash.Sum([]byte{})
}
func (s *cfgState) xorItemHash(itemHash []byte) {
	for i := 0; i < 32; i++ {
		s.checksum[i] = s.checksum[i] ^ itemHash[i]
	}
}

func (s *cfgState) Add(files map[string][]byte) State {
	newState := s.clone()

	for name, file := range files {
		pathParts := deleteEmpty(strings.Split(name, "/"))
		isHidden := false
		for _, part := range pathParts {
			if part[0] == '.' {
				isHidden = true
				break
			}
		}
		oldContent := newState.rootDir.replaceFile(pathParts, file, true)
		if !isHidden {
			pathString := name
			oldHash := checksum(pathString, isHidden, oldContent)
			newHash := checksum(pathString, isHidden, file)
			newState.xorItemHash(oldHash)
			newState.xorItemHash(newHash)
		}
		newState.size += uint(len(file) - len(oldContent))
	}
	return newState
}

func EmptyState() *cfgState {
	return &cfgState{
		checksum: make([]byte, 32),
		size:     0,
		rootDir:  &cfgDir{},
	}
}

func (s *cfgState) Checksum() string {
	return hex.EncodeToString(s.checksum)
}

var _ State = &cfgState{}
