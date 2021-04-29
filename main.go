package main

import (
	"Heterogenous_SRM/migrator"
	"Heterogenous_SRM/watcher"
	"sync"
	"time"
)

func main() {
	// declare the needed structures
	done := make(chan bool)
	fileAccessMap := make(map[string][]time.Time)
	fileAge := make(map[string]time.Time)
	storagePolicy := make(map[string]string)
	var mutex = &sync.Mutex{}

	// start migrator and watcher
	go migrator.StartMigrator(storagePolicy, fileAccessMap, fileAge, mutex)
	go watcher.WatcherFunc("./test", fileAccessMap, fileAge, mutex)

	<-done
}
