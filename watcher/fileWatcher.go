package watcher

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatcherFunc(path string, fileAccess map[string][]time.Time) {

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	dir := getDirName(path)

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				eventName := event.Op.String()
				dirfileName := event.Name
				fileName := getFileName(dirfileName)

				if fileName == dir {
					continue
				}

				capTimeStampsForOneMonth(fileAccess, fileName)

				if eventName == "OPEN" || eventName == "WRITE" || eventName == "CREATE" {
					ts := time.Now()
					if fileAccess[fileName] != nil {
						fileAccess[fileName] = append(fileAccess[fileName], ts)
					} else {
						fileAccess[fileName] = []time.Time{ts}
					}
				}

				// watch for errors
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
