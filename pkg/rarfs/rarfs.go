package rarfs

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/nwaples/rardecode"
	"github.com/psanford/memfs"
)

func New(r io.Reader) (fs.FS, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	arch, err := rardecode.NewReader(bytes.NewReader(data), "")
	if err != nil {
		return nil, err
	}

	fsys := memfs.New()

	for {
		header, err := arch.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		if header.IsDir {
			err := fsys.MkdirAll(header.Name, header.Mode())
			if err != nil {
				return nil, err
			}
		} else {
			dir := filepath.Dir(header.Name)

			f, err := fsys.Open(dir)
			if err != nil {
				// the directory does not exist so we must create it first
				err = fsys.MkdirAll(dir, header.Mode())
				if err != nil {
					return nil, err
				}
			} else {
				// the directory exists
				f.Close()
			}

			fdata, err := io.ReadAll(arch)
			if err != nil {
				return nil, err
			}

			err = fsys.WriteFile(header.Name, fdata, header.Mode())
			if err != nil {
				return nil, err
			}
		}
	}

	return fsys, nil
}
