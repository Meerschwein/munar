package main

import (
	"log"
	"os"
	"strings"

	"github.com/meerschwein/unar/pkg/archives"
	"github.com/meerschwein/unar/pkg/util"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) != 2 {
		log.Fatalf("Usage: unpack <filename>")
	}

	filename := os.Args[1]

	for _, a := range archives.Archives {
		if strings.HasSuffix(filename, a.Suffix) {
			dstPath := util.CleanPath(filename, a.Suffix)
			util.CheckDestinationDir(dstPath)

			file, err := os.Open(filename)
			if err != nil {
				log.Fatalf("failed to open file %s: %v", filename, err)
			}

			fs, err := a.ToFs(file)
			if err != nil {
				log.Fatalf("failed to unpack file %s: %v", filename, err)
			}

			err = util.CopyFS(dstPath, fs)
			if err != nil {
				log.Fatalf("failed to unpack file %s: %v", filename, err)
			}

			file.Close()
			return
		}
	}

	log.Fatalf("unknown file extension")
}
