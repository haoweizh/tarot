package util

import (
	"time"
)


func GetNow() time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		return time.Now().In(location)
	}
	return time.Now()
}
