package watcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatcherFunc(path string, fileAccess map[string][]time.Time, fileAge map[string]time.Time, mutex *sync.Mutex) {

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

				capTimeStampsForOneMonth(fileAccess, fileName, mutex)

				if eventName == "OPEN" || eventName == "WRITE" || eventName == "CREATE" {
					ts := time.Now()
					if fileAccess[fileName] != nil {
						mutex.Lock()
						fileAccess[fileName] = append(fileAccess[fileName], ts)
						mutex.Unlock()
					} else {
						mutex.Lock()
						fileAccess[fileName] = []time.Time{ts}
						fileAge[fileName] = ts
						mutex.Unlock()
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
