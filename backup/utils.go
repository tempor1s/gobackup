package backup

import (
	"fmt"
	"io/ioutil"
	"log"
)

// Get the size of a directory in the format of a string.
func getDirSizeStr(dirName string) string {
	size := dirSize(dirName)
	sizeStr := fmt.Sprintf("%.2f MB", size)

	if size > 1000 {
		size = size / 1000
		sizeStr = fmt.Sprintf("%.2f GB", size)
	}

	return sizeStr
}

// fileSize is a var we keep because our other method is recursive, and it just makes life easier
var fileSize float64

// Get the size of a directory in MB, returns a float
func dirSize(dirName string) float64 {
	allFiles, err := ioutil.ReadDir(dirName)
	checkIfError(err)

	for _, file := range allFiles {
		if file.IsDir() {
			dirSize(dirName + "/" + file.Name())
		}
		fileSize += float64(file.Size()) / 1000000.0
	}

	return fileSize
}

// Just a simple check to see if there is an error, because I was re-writing this code too much
func checkIfError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
