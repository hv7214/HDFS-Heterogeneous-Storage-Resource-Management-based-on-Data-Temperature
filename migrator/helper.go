package migrator

import (
	"Heterogenous_SRM/database"
	"time"
)

var dayMilliseconds int64 = int64(86400000)
var weekMilliseconds int64 = int64(604800000)
var monthMilliseconds int64 = int64(2592000000)
var yearMilliseconds int64 = int64(31536000000)

func getCountMetrics(accessTimes []time.Time) (int, int, int) {
	count_m := len(accessTimes)
	count_w := 0
	count_d := 0
	boolWeek := false
	boolDay := false

	for ind := 0; ind < len(accessTimes); ind++ {
		if time.Since(accessTimes[ind]).Milliseconds() <= weekMilliseconds && !boolWeek {
			count_w = len(accessTimes) - ind + 1
			boolWeek = true
		}
		if time.Since(accessTimes[ind]).Milliseconds() <= dayMilliseconds && !boolDay {
			count_d = len(accessTimes) - ind + 1
			boolDay = true
		}
	}

	return count_d, count_w, count_m
}

func getTemperature(count_d int, count_m int, count_w int, age time.Time) string {
	temperature := "N/A"
	ageInMilliseconds := time.Since(age).Milliseconds()
	if ageInMilliseconds < weekMilliseconds && count_d > 30 {
		temperature = "SUMMER"
	} else if ageInMilliseconds > weekMilliseconds && ageInMilliseconds < monthMilliseconds && count_d > 15 && count_w > 30 {
		temperature = "FALL"
	} else if ageInMilliseconds > monthMilliseconds && ageInMilliseconds < 3*monthMilliseconds && count_w == 0 && count_m > 0 {
		temperature = "WINTER"
	} else if ageInMilliseconds > 3*monthMilliseconds && ageInMilliseconds < yearMilliseconds && count_m == 0 {
		temperature = "FROZEN"
	}

	return temperature
}

func capTimeStampsForOneMonth(fileAccess map[string][]time.Time) {
	var ind int
	monthMilliseconds := int64(2592000000)

	for filename := range fileAccess {
		for ind = 0; ind < len(fileAccess[filename]); ind++ {
			if time.Since(fileAccess[filename][ind]).Milliseconds() < monthMilliseconds {
				break
			}
		}
		fileAccess[filename] = fileAccess[filename][ind:]
		if ind > 0 {
			database.UpdateAccess(filename, fileAccess[filename])
		}
	}
}
