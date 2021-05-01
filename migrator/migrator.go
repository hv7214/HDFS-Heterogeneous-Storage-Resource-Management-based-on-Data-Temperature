package migrator

import (
	"fmt"
	"math"
	"sync"
	"time"
)

var Temp_to_storage_policy map[string]string = map[string]string{
	"SUMMER": "ALL_SSD",
	"FALL":   "ONE_SSD",
	"WINTER": "WARM",
	"FROZEN": "COLD",
	"N/A":    "LAZY PERSIST",
}

// testing
// var migrator_run_interval time.Duration = time.Second

// production
var migrator_run_interval time.Duration = time.Hour

func StartMigrator(storagePolicy map[string]string, fileAccess map[string][]time.Time, fileAge map[string]time.Time, mutex *sync.Mutex) {
	totalAccessInADay := 1
	ticker := time.NewTicker(migrator_run_interval)

	for {
		select {
		case <-ticker.C:
			fmt.Println("here")
			// lock the mutex
			mutex.Lock()
			for filename, accessTimes := range fileAccess {
				// get count metrics
				count_d, count_w, count_m := getCountMetrics(accessTimes)
				// add count_d value to get avg in the end
				totalAccessInADay += count_d
				// get new temperature of the file
				temperature := getTemperature(count_d, count_w, count_m, fileAge[filename])
				// get the new storage policy depending on the temperature
				newStoragePolicy := Temp_to_storage_policy[temperature]
				// if new storage policy is not same as before, invoke the mover
				if newStoragePolicy != storagePolicy[filename] {
					fmt.Println("Invoking mover: " + filename + " storage policy changed from " +
						storagePolicy[filename] + " to " + newStoragePolicy)
				}
				// update the storage policy
				storagePolicy[filename] = newStoragePolicy
			}
			// unlock the mutex
			mutex.Unlock()
			// update the migrator run interval time by 24hrs/totalAccessInADay
			timeTakenForOneAccess := fmt.Sprintf("%f", math.Min(float64((24*3600)/totalAccessInADay), 60)) + "s"
			// update the migrator run interval value(it should not go beyond one minute)
			if (24*3600)/totalAccessInADay > 60 {
				migrator_run_interval, _ = time.ParseDuration(timeTakenForOneAccess)
			} else {
				migrator_run_interval = time.Minute
			}
			ticker = time.NewTicker(migrator_run_interval)
		}
	}
}
