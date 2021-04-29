package migrator

import "time"

func startMigrator(storagePolicy map[string]string, fileAccess map[string][]time.Time, fileAge map[string]time.Time) {
	for filename, accessTimes := range fileAccess {
		count_d, count_w, count_m := getCountMetrics(accessTimes)
		temperature := getTemperature(count_d, count_w, count_m, fileAge[filename])
	}
}
