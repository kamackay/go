package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"path/filepath"
)

func dirSize(path string) (int64, int64, error) {
	var size int64
	var fileCount int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
			fileCount++
		}
		return err
	})
	return size, fileCount, err
}

func main() {
	var args []string
	for _, s := range os.Args[1:] {
		args = append(args, strings.ToLower(s))
	}
	log := func(s ...interface{}) { fmt.Println(s...) } // Make code shorter
	startTime := time.Now()
	if len(args) == 0 {
		log("No Command Given")
	} else {
		switch args[0] {
		case "filesize":
			if len(args) < 2 {
				log("No Path Given")
				break
			}
			path := args[1]
			log("Process size of", path)
			size, files, err := dirSize(path)
			if err != nil {
				panic(err)
			}
			log("\t", size, "bytes\t\t", files, "files")

			var fSize float64 = float64(size)
			var unit uint8 = 0
			for fSize > 1024 {
				fSize /= 1024
				unit++
			}
			log("\t", fSize, []string{"bytes", "KB", "MB", "GB", "TB"}[unit])
			break
		default:
			log("Unrecognized Command:", os.Args[1]) // Print the original command
			break
		}
	}
	defer log("Finished in", time.Now().Sub(startTime).String())
}
