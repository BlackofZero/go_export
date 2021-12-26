package utils

import (
	"time"
)

func GetNowDayStr() string {
	return time.Now().Format("2006-01-02")
}
