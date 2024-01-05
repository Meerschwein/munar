package archives

import (
	"archive/zip"
	"compress/gzip"
	"io/fs"
	"os"

	"github.com/bodgit/sevenzip"
	"github.com/nlepage/go-tarfs"
	"github.com/ulikunitz/xz"
)

type Archive struct {
	Suffix string
	ToFs   func(src *os.File) (fs.FS, error)
}

var Archives = []Archive{
	{".tar", tarFs},
	{".tar.gz", tarGzFs},
	{".tar.xz", tarXzFs},
	{".zip", zipFs},
	{".7z", sevenZipFs},
}

func tarFs(src *os.File) (fs.FS, error) {
	return tarfs.New(src)
}

func tarGzFs(src *os.File) (fs.FS, error) {
	gzipReader, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}
	return tarfs.New(gzipReader)
}

func tarXzFs(src *os.File) (fs.FS, error) {
	xzReader, err := xz.NewReader(src)
	if err != nil {
		return nil, err
	}
	return tarfs.New(xzReader)
}

func zipFs(src *os.File) (fs.FS, error) {
	info, err := src.Stat()
	if err != nil {
		return nil, err
	}
	return zip.NewReader(src, info.Size())
}

func sevenZipFs(src *os.File) (fs.FS, error) {
	info, err := src.Stat()
	if err != nil {
		return nil, err
	}
	return sevenzip.NewReader(src, info.Size())
}
