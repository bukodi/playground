package config

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"embed"
	_ "embed"
	"errors"
	"io"
	"os"
	"path/filepath"
)

//go:embed testdata/cfgdir/.meta
var metaDir embed.FS

func ImportDir(basePath string) (State, error) {

	files := make(map[string][]byte)

	var readDir func(path string) error
	readDir = func(path string) error {
		dirEntries, err := os.ReadDir(filepath.Join(basePath, path))
		if err != nil {
			return err
		}
		for _, de := range dirEntries {
			entryPath := filepath.Join(path, de.Name())
			if de.IsDir() {
				if err := readDir(entryPath); err != nil {
					return err
				}
			} else {
				if fileContent, err := os.ReadFile(filepath.Join(basePath, entryPath)); err != nil {
					return err
				} else if len(fileContent) > 0 {
					files[entryPath] = fileContent
				}
			}
		}
		return nil
	}
	if err := readDir(""); err != nil {
		return nil, err
	}

	return EmptyState().Add(files), nil
}

func ExportTGZ(s State, out io.WriteCloser, addMeta bool) (retErr error) {

	zw := gzip.NewWriter(out)
	tw := tar.NewWriter(zw)
	defer func() {
		retErr = errors.Join(retErr, tw.Close(), zw.Close())
	}()

	files := s.List("")
	for path, content := range files {
		// generate tar header
		if err := tw.WriteHeader(&tar.Header{
			Name:     path,
			Typeflag: tar.TypeReg,
			Size:     int64(len(content)),
		}); err != nil {
			return err
		}

		if _, err := tw.Write(content); err != nil {
			return err
		}
	}

	if addMeta {
		entries, err := metaDir.ReadDir(".meta")
		if err != nil {
			return err
		}
		for _, e := range entries {
			e.IsDir()
		}
	}

	return nil
}

func ExportZip(s State, out io.WriteCloser, addMeta bool) (retErr error) {
	zw := zip.NewWriter(out)
	defer func() {
		retErr = errors.Join(retErr, zw.Close())
	}()

	files := s.List("")
	for path, content := range files {
		zew, err := zw.Create(path)
		if err != nil {
			return err
		}

		if _, err := zew.Write(content); err != nil {
			return err
		}
	}

	return nil
}
