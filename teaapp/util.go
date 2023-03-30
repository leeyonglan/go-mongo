package teaapp

import "time"

func GetLocalTime(zoneOffset int) time.Time {
	gmt := time.Now().UTC()
	local := gmt.Add(time.Duration(zoneOffset * int(time.Hour)))
	return local
}

func isInTime(currentTime time.Time) bool {
	hour := currentTime.UTC().Hour()
	if (hour >= 12 && hour <= 14) || (hour >= 18 && hour <= 20) || (hour >= 22 && hour <= 23) {
		return true
	}
	return false
}

func isForbiddenTime(currentTime time.Time) bool {
	hour := currentTime.UTC().Hour()
	if hour >= 0 && hour <= 7 {
		return true
	}
	return false
}
