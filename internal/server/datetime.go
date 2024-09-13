package server

import "time"

func dateTimeNow() string {
	tz := time.FixedZone("Asia/Tokyo", 9*3600)
	return time.Now().In(tz).Format(time.DateTime)
}
