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

func FilterByte(content string, old byte) string {
	bytes := []byte(content)
	newBytes := make([]byte, 0)
	for _, value := range bytes {
		if value != old {
			newBytes = append(newBytes, value)
		}
	}
	return string(newBytes)
}
