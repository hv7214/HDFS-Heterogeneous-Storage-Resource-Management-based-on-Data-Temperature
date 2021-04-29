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
	"N/A":    "HOT",
}

var migrator_run_interval time.Duration = time.Hour

func StartMigrator(storagePolicy map[string]string, fileAccess map[string][]time.Time, fileAge map[string]time.Time, mutex *sync.Mutex) {
	totalAccessInADay := 0

	for _ = range time.Tick(migrator_run_interval) {
		// lock the mutex
		mutex.Lock()
		for filename, accessTimes := range fileAccess {
			// get count metrics
			count_d, count_w, count_m := getCountMetrics(accessTimes)
			// add count_d value to get avg in the end
			totalAccessInADay += count_d
			// get new temperature of the file
			temperature := getTemperature(count_d, count_w, count_m, fileAge[filename])
			// update the storage policy depending on the temperature
			storagePolicy[filename] = Temp_to_storage_policy[temperature]
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
	}
}
