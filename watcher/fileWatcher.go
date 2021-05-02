package watcher

import (
	"fmt"
	"time"

	"Heterogenous_SRM/database"

	"github.com/fsnotify/fsnotify"
)

func WatcherFunc(path string) {

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	dir := getDirName(path)

	done := make(chan bool)

	go func() {
		for {
			pass := 0
			prevfile := ""
			select {
			// watch for events
			case event := <-watcher.Events:
				eventName := event.Op.String()
				dirfileName := event.Name
				fileName := getFileName(dirfileName)

				if fileName == dir {
					continue
				}

				if pass > 0 && fileName == prevfile {
					pass = pass - 1
					continue
				}

				if eventName == "WRITE" {
					prevfile = fileName
					pass = 3
				}

				// capTimeStampsForOneMonth(fileAccess, fileName, mutex)

				if eventName == "OPEN" || eventName == "WRITE" || eventName == "CREATE" {
					ts := time.Now()
					flag, data := database.CheckExists(fileName)
					if flag == true {
						data = append(data, ts)
						database.UpdateAccess(fileName, data)
					} else {
						database.InsertAccessAndAge(fileName, ts, "HOT")
						// storagePolicy[fileName] = "HOT"
					}
				}

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(path); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
