package watcher

import (
	"strings"
	"time"
)

func getDirName(path string) string {
	ind := strings.LastIndex(path, "/")
	dir := path[ind+1:]
	return dir
}

func getFileName(dirfileName string) string {
	ind := strings.Index(dirfileName, "/")
	fileName := dirfileName[ind+1:]
	return fileName
}

func capTimeStampsForOneMonth(fileAccess map[string][]time.Time, filename string) {
	var ind int
	monthMilliseconds := int64(2629800000)

	for ; ind < len(fileAccess[filename]); ind++ {
		if time.Since(fileAccess[filename][ind]).Milliseconds() < monthMilliseconds {
			break
		}
	}

	fileAccess[filename] = fileAccess[filename][ind:]
}
