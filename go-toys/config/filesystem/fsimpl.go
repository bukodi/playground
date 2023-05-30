package filesystem

import (
	"crypto"
	"path/filepath"
)

var separator string = string(filepath.Separator)
var hash crypto.Hash = crypto.SHA256

type fsCfgState struct {
}

type cfgEntry struct {
	parent  *cfgDir
	subPath string
}

func (e *cfgEntry) path() string {
	if e.parent == nil {
		return e.subPath
	} else {
		parentPath := e.parent.path()
		return parentPath + separator + e.subPath
	}
}

type cfgDir struct {
	cfgEntry
	entries cfgEntry
}

type cfgValue struct {
	cfgEntry
	value []byte
}

func (d cfgDir) checksum() string {

	panic("should not call directly")
}

func (e cfgValue) checksum() []byte {
	h := hash.New()
	h.Write([]byte(e.path()))
	h.Write(e.value)
	return h.Sum(nil)
}
