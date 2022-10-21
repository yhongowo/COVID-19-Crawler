package util

import "time"

// 得到当天0点的时间戳
func TodayBeginTime() int64 {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return today.Unix()
}
