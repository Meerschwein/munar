package util

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CleanPath returns the base name of path without its suffix.
func CleanPath(path, suffix string) string {
	return filepath.Base(strings.TrimSuffix(path, suffix))
}

// CheckDestinationDir bail out if the destination directory already exists.
func CheckDestinationDir(dstDir string) {
	file, err := os.OpenFile(dstDir, os.O_RDONLY, 0)
	if err == nil {
		file.Close()
		log.Fatalf("directory %s already exists", dstDir)
	}
}

// shamelessly stolen from https://github.com/golang/go/issues/62484
func CopyFS(dir string, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, _err error) error {
		targ := filepath.Join(dir, filepath.FromSlash(path))
		if d.IsDir() {
			if err := os.MkdirAll(targ, 0o777); err != nil {
				return err
			}
			return nil
		}
		r, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()
		info, err := r.Stat()
		if err != nil {
			return err
		}
		w, err := os.OpenFile(targ, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666|info.Mode()&0o777)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			w.Close()
			return fmt.Errorf("copying %s: %v", path, err)
		}
		if err := w.Close(); err != nil {
			return err
		}
		return nil
	})
}
