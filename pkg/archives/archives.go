package archives

import (
	"io"
	"io/fs"

	"github.com/bodgit/sevenzip"
	"github.com/josharian/txtarfs"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zip"
	"github.com/meerschwein/unar/internal/rarfs"
	"github.com/nlepage/go-tarfs"
	"github.com/ulikunitz/xz"
	"golang.org/x/tools/txtar"
)

type Reader interface {
	io.Reader
	io.ReaderAt
}

type ArchiveFsFn func(r Reader, size int64) (fs.FS, error)

var SuffixArchives = map[string]ArchiveFsFn{
	".7z":     sevenZipFs,
	".epub":   zipFs,
	".odt":    zipFs,
	".rar":    rarFs,
	".tar.gz": tarGzFs,
	".tar.xz": tarXzFs,
	".tar":    tarFs,
	".tgz":    tarGzFs,
	".txtar":  txtarFs,
	".xpi":    zipFs,
	".zip":    zipFs,
}

var FormatArchives = map[string]ArchiveFsFn{
	"7zip":  sevenZipFs,
	"rar":   rarFs,
	"tar":   tarFs,
	"targz": tarGzFs,
	"tarxz": tarXzFs,
	"txtar": txtarFs,
	"zip":   zipFs,
}

func tarFs(r Reader, _ int64) (fs.FS, error) {
	return tarfs.New(r)
}

func tarGzFs(r Reader, _ int64) (fs.FS, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return tarfs.New(gzipReader)
}

func tarXzFs(r Reader, _ int64) (fs.FS, error) {
	xzReader, err := xz.NewReader(r)
	if err != nil {
		return nil, err
	}
	return tarfs.New(xzReader)
}

func zipFs(r Reader, size int64) (fs.FS, error) {
	return zip.NewReader(r, size)
}

func sevenZipFs(r Reader, size int64) (fs.FS, error) {
	return sevenzip.NewReader(r, size)
}

func txtarFs(r Reader, _ int64) (fs.FS, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return txtarfs.As(txtar.Parse(content)), nil
}

func rarFs(r Reader, _ int64) (fs.FS, error) {
	return rarfs.New(r)
}
