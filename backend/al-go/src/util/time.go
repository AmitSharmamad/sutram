package util

import "time"

func UnixToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func ToClickhouse(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
