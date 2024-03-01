package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/meerschwein/unar/pkg/archives"
	"github.com/meerschwein/unar/pkg/util"
	"golang.org/x/exp/maps"
)

func init() {
	log.SetFlags(0)

	fmts := maps.Keys(archives.FormatArchives)
	slices.Sort(fmts)
	flag.StringVar(&format, "f", "", "archive format of the file\npossible values: "+strings.Join(fmts, ", "))

	flag.Usage = func() {
		log.Println("Usage: unar [options] archive")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename = flag.Arg(0)
}

var (
	format   string
	filename string
)

func main() {
	var archFsFn archives.ArchiveFsFn
	var dstPath string

	if format != "" {
		found := false
		archFsFn, found = archives.FormatArchives[format]
		if !found {
			log.Fatal("unknown archive format ", format)
		}
		dstPath = filepath.Base(strings.TrimSuffix(filename, filepath.Ext(filename)))
	} else {
		for suffix, fn := range archives.SuffixArchives {
			if strings.HasSuffix(filename, suffix) {
				archFsFn = fn
				dstPath = util.CleanPath(filename, suffix)
				goto found
			}
		}

		log.Println("unknown file extension")
		flag.Usage()
		os.Exit(1)

	found:
	}

	_, err := os.Open(dstPath)
	if err == nil {
		log.Fatal(`directory "`, dstPath, `" already exists`)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fs, err := archFsFn(file)
	if err != nil {
		log.Fatal(err)
	}

	err = util.CopyFS(dstPath, fs)
	if err != nil {
		log.Fatal(err)
	}
}
