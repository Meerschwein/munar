package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fatalf("Usage: unpack <filename>")
	}

	filename := flag.Arg(0)

	switch {
	case strings.HasSuffix(filename, ".zip"):
		unpackZip(filename)
	case strings.HasSuffix(filename, ".tar.gz"):
		unpackTarGz(filename)
	default:
		fatalf("unknown file extension")
	}
}

func unpackTarGz(filename string) {
	dstDir := strings.TrimSuffix(filename, ".tar.gz") // remove '.tar.gz'
	dstDir = filepath.Base(dstDir)                    // unpack it to the current directory

	if _, err := os.OpenFile(dstDir, os.O_RDONLY, 0); err == nil {
		fatalf("directory %s already exists", dstDir)
	}

	file, err := os.Open(filename)
	if err != nil {
		fatalf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		fatalf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// skip pax header
	_, err = tarReader.Next()
	if err != nil {
		fatalf("failed to read tar header: %v", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			fatalf("failed to read tar header: %v", err)
		}

		fmt.Printf("unpacking %s", header.Name)

		dstFilename := filepath.Join(dstDir, header.Name)

		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(dstFilename, 0o755); err != nil {
				fatalf("failed to create directory %s: %v", dstFilename, err)
			}
		} else {
			dstParentDir := filepath.Dir(dstFilename)

			if err := os.MkdirAll(dstParentDir, 0o755); err != nil {
				fatalf("failed to create directory %s: %v", dstParentDir, err)
			}

			dstFile, err := os.OpenFile(dstFilename, os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				fatalf("failed to create file %s: %v", dstFilename, err)
			}

			_, err = io.Copy(dstFile, tarReader)
			if err != nil {
				fatalf("failed to write to %s: %v", dstFilename, err)
			}

			dstFile.Close()
		}
	}
}

func unpackZip(filename string) {
	dstDir := strings.TrimSuffix(filename, ".zip") // remove '.zip'
	dstDir = filepath.Base(dstDir)                 // unpack it to the current directory

	if _, err := os.OpenFile(dstDir, os.O_RDONLY, 0); err == nil {
		fatalf("directory %s already exists", dstDir)
	}

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		fatalf("failed to open zip file %s: %v", filename, err)
	}
	defer zipReader.Close()

	if err := os.Mkdir(dstDir, 0o755); err != nil {
		fatalf("failed to create destination directory %s: %v", dstDir, err)
	}

	for _, fileInZip := range zipReader.File {
		fmt.Printf("unpacking %s", fileInZip.Name)
		srcFile, err := fileInZip.Open()
		if err != nil {
			fatalf("failed to open file %s in zip file %s: %v", fileInZip.Name, filename, err)
		}

		dstFilename := filepath.Join(dstDir, fileInZip.Name)

		if fileInZip.FileHeader.FileInfo().IsDir() {
			if err := os.MkdirAll(dstFilename, 0o755); err != nil {
				fatalf("failed to create directory %s: %v", dstFilename, err)
			}
		} else {
			dstParentDir := filepath.Dir(dstFilename)

			if err := os.MkdirAll(dstParentDir, 0o755); err != nil {
				fatalf("failed to create directory %s: %v", dstParentDir, err)
			}

			dstFile, err := os.OpenFile(dstFilename, os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				fatalf("failed to create file %s: %v", dstFilename, err)
			}

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				fatalf("failed to write to %s: %v", dstFilename, err)
			}

			dstFile.Close()
		}

		srcFile.Close()
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
