package common

import (
	"time"
)

func ZeroAclock() time.Time {
	str := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", str, time.Local)
	return t.AddDate(0, 0, 1)
}
