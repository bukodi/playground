package config_test

import (
	"crypto/sha256"
	"errors"
	"github.com/bukodi/playground/config"
	"os"
	"path/filepath"
	"testing"
)

func TestEmptyState(t *testing.T) {
	s0 := config.EmptyState()
	dump(t, s0, "s0")
	s1 := s0.Add(map[string][]byte{
		"foo/bar/data.json": []byte{1, 2, 3},
	})
	dump(t, s1, "s1")
	s2 := s1.Add(map[string][]byte{
		"foo/bar/data.json": nil,
	})
	dump(t, s2, "s2")
}

func dump(t *testing.T, s config.State, name string) {
	t.Logf("%s : %s, %d", name, s.Checksum(), s.Size())
	files := s.List("")
	if len(files) > 0 {
		for n, f := range files {
			t.Logf("    %s : %d bytes", n, len(f))
		}
	} else {
		t.Logf("    no files")
	}
}

func TestFSImport(t *testing.T) {
	s, err := config.ImportDir(filepath.Join("testdata", "cfgdir"))
	if err != nil {
		t.Errorf("%+v", err)
	}
	dump(t, s, "imported")

	tgzFile, err := os.Create("/tmp/test.tgz")
	if err != nil {
		t.Errorf("%+v", err)
	}
	err = config.ExportTGZ(s, tgzFile, false)
	if err != nil {
		t.Errorf("%+v", err)
	}

	zipFile, err := os.Create("/tmp/test.zip")
	if err != nil {
		t.Errorf("%+v", err)
	}
	err = config.ExportZip(s, zipFile, false)
	if err != nil {
		t.Errorf("%+v", err)
	}
}

func TestFSExport(t *testing.T) {
	s, err := config.ImportDir(filepath.Join("testdata", "cfgdir"))
	if err != nil {
		t.Errorf("%+v", err)
	}
	dump(t, s, "imported")
}

func TestHashSum(t *testing.T) {

	{
		h := sha256.New()
		h.Write([]byte{1, 2, 3})
		h.Write([]byte{4, 5, 6})
		t.Logf("%x", h.Sum(nil))
	}
	{
		h := sha256.New()
		h.Write([]byte{1, 2, 3, 4, 5, 6})
		t.Logf("%x", h.Sum(nil))
	}
}

func TestDeferErr(t *testing.T) {
	fn := func() (retErr error) {
		defer func() { retErr = errors.Join(retErr, errors.New("Error in defer")) }()
		return errors.New("Error in return")
	}
	err := fn()
	t.Logf("%+v", err)
}
