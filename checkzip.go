package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func openzip(path, brokenDir string) (result bool) {
	defer func() {
		recover()
	}()

	reader, err := zip.OpenReader(path)
	// after recovery from panic
	defer reader.Close()
	if err != nil {
		fmt.Printf("err = %+v\n", err)
		return false
	}
	return true
}

func main() {
	var dir string = ""

	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		dir = args[0]
	} else {
		log.Fatal("dir not specified")
	}

	brokenDir := filepath.Join(dir, "broken")
	os.Mkdir(brokenDir, 0755)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == dir || strings.HasPrefix(path, brokenDir) || filepath.Ext(path) != ".zip" {
			return nil
		}

		fmt.Printf("start checking: %+v\n", path)
		result := openzip(path, brokenDir)
		if result {
			fmt.Printf("%+v is valid.\n", path)
		} else {
			fmt.Printf("%+v is broken, moved under broken/ directory.\n", path)
			os.Rename(path, filepath.Join(brokenDir, filepath.Base(path)))
		}
		return nil
	})
}
