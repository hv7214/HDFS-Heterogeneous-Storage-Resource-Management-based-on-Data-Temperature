package main

import (
	"Heterogenous_SRM/watcher"
	"time"
)

func main() {
	done := make(chan bool)
	fileAccessMap := make(map[string][]time.Time)
	fileAge := make(map[string]time.Time)
	go watcher.WatcherFunc("./test", fileAccessMap, fileAge)
	<-done
}
