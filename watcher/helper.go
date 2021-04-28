package watcher

import "strings"

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
