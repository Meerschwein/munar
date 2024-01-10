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
		log.Fatal("Usage: unpack <filename>")
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, a := range archives.Archives {
		if strings.HasSuffix(filename, a.Suffix) {
			dstPath := util.CleanPath(filename, a.Suffix)
			util.CheckDestinationDir(dstPath)

			fs, err := a.ToFs(file)
			if err != nil {
				log.Fatal(err)
			}

			err = util.CopyFS(dstPath, fs)
			if err != nil {
				log.Fatal(err)
			}

			return
		}
	}

	log.Fatal("unknown file extension")
}
