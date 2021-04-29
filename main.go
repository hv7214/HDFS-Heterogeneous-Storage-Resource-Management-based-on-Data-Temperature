package main

import (
	"Heterogenous_SRM/watcher"
	"time"
	"Heterogenous_SRM/migrator"
	"sync"
)

func main() {
	done := make(chan bool)
	fileAccessMap := make(map[string][]time.Time)
	fileAge := make(map[string]time.Time)
	storagePolicy := make(map[string]string)
	var mutex = &sync.Mutex{}
	go migrator.StartMigrator(storagePolicy, fileAccessMap, fileAge, mutex)
	go watcher.WatcherFunc("./test", fileAccessMap, fileAge, mutex)
	<-done
}
