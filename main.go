package main

import (
	"Heterogenous_SRM/database"
	"Heterogenous_SRM/migrator"
	"Heterogenous_SRM/watcher"
)

func main() {
	// declare the needed structures
	database.ConnectToDb()
	done := make(chan bool)

	// start migrator and watcher
	go migrator.StartMigrator()
	go watcher.WatcherFunc("./test")

	<-done
}
