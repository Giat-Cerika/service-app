package utils

import "time"

func FormatDate(t time.Time) string {
	return t.Format("02-01-2006")
}

func FormatTime(t time.Time) string {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	return t.In(locJakarta).Format("03:04 PM")
}
