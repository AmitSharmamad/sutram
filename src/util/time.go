package util

import "time"

func UnixToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func ToClickhouse(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func toUnix(datetimeStr string, fromFormat string) (int64, error) {
	t, err := time.Parse(fromFormat, datetimeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func fromUnix(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2023-03-26 00:30:00")
}
