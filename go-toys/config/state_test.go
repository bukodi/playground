package config_test

import (
	"github.com/bukodi/playground/config"
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
			t.Logf("    %s : %+v", n, f)
		}
	} else {
		t.Logf("    no files")
	}
}
