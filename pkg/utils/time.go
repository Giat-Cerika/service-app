package utils

import "time"

func FormatDate(t time.Time) string {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	return t.In(locJakarta).Format("02-01-2006")
}

func FormatTime(t time.Time) string {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	return t.In(locJakarta).Format("03:04 PM")
}

func FormatDateTime(t *time.Time) string {
	if t == nil {
		return ""
	}

	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	return t.In(locJakarta).Format("02-01-2006 03:04 PM")
}
func FormatOnlyDate(t time.Time) string {
	return t.Format("02-01-2006")
}

func FormatLogDate(dateStr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return dateStr // fallback
	}
	return t.Format("2006-01-02")
}
