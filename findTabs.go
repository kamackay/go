package main

import (
	"io/ioutil"
	"fmt"
	"os"
	"path"
)

func getFiles(dir string) []string {
	var files []string
	fileList, _ := ioutil.ReadDir(dir)
	for _, f := range fileList {
		if !f.IsDir() {
			files = append(files, path.Join(dir, f.Name()))
		} else {
			for _, subFile := range getFiles(path.Join(dir, f.Name())) {
				files = append(files, subFile)
			}
		}
	}
	return files
}

func main() {
	log := func(s ...interface{}) { fmt.Println(s...) }
	args := os.Args[1:]
	dir, _ := os.Getwd()
	if len(args) > 0 {
		dir = args[0]
	}

	log("Searching for tabs in", dir)
	for _, f := range getFiles(dir) {
		log(f)
	}
}
