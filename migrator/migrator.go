package migrator

import (
	"time"
	"fmt"
	"sync"
)

var Temp_to_storage_policy map[string]string = map[string]string{
	"SUMMER" : "ALL_SSD",
	"FALL" : "ONE_SSD",
	"WINTER" : "WARM",
	"FROZEN" : "COLD",
	"N/A" : "HOT",
}

var migrator_run_interval time.Duration = time.Minute

func StartMigrator(storagePolicy map[string]string, fileAccess map[string][]time.Time, fileAge map[string]time.Time, mutex *sync.Mutex) {
	for _= range time.Tick(migrator_run_interval){
		mutex.Lock()
		for filename, accessTimes := range fileAccess {
			count_d, count_w, count_m := getCountMetrics(accessTimes)
			fmt.Println(filename)
			fmt.Println(count_d)
			temperature := getTemperature(count_d, count_w, count_m, fileAge[filename])
			storagePolicy[filename] = Temp_to_storage_policy[temperature]
		}
		mutex.Unlock()
		fmt.Println(storagePolicy)
	}
}
