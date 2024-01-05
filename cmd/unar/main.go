package main

import (
	"archive/zip"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		logger.Fatalf("Usage: unpack <filename>")
	}

	filename := flag.Arg(0)

	switch {
	case strings.HasSuffix(filename, ".zip"):
		unpackZip(filename)
	default:
		logger.Fatalf("unknown file extension")
	}
}

func unpackZip(filename string) {
	dstDir := strings.TrimSuffix(filename, ".zip") // remove '.zip'
	dstDir = filepath.Base(dstDir)                 // unpack it to the current directory

	if file, err := os.OpenFile(dstDir, os.O_RDONLY, 0); err == nil {
		file.Close()
		logger.Fatalf("directory %s already exists", dstDir)
	}

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		logger.Fatalf("failed to open zip file %s: %v", filename, err)
	}
	defer zipReader.Close()

	if err := os.Mkdir(dstDir, 0o755); err != nil {
		logger.Fatalf("failed to create destination directory %s: %v", dstDir, err)
	}

	for _, fileInZip := range zipReader.File {
		logger.Printf("unpacking %s", fileInZip.Name)
		srcFile, err := fileInZip.Open()
		if err != nil {
			logger.Fatalf("failed to open file %s in zip file %s: %v", fileInZip.Name, filename, err)
		}

		dstFilename := filepath.Join(dstDir, fileInZip.Name)

		if fileInZip.FileHeader.FileInfo().IsDir() {
			if err := os.MkdirAll(dstFilename, 0o755); err != nil {
				logger.Fatalf("failed to create directory %s: %v", dstFilename, err)
			}
		} else {
			dstParentDir := filepath.Dir(dstFilename)

			if err := os.MkdirAll(dstParentDir, 0o755); err != nil {
				logger.Fatalf("failed to create directory %s: %v", dstParentDir, err)
			}

			dstFile, err := os.OpenFile(dstFilename, os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				logger.Fatalf("failed to create file %s: %v", dstFilename, err)
			}

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				logger.Fatalf("failed to write to %s: %v", dstFilename, err)
			}

			dstFile.Close()
		}

		srcFile.Close()
	}
}
