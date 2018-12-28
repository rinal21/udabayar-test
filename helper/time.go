package helper

import (
	"fmt"
	"time"
)

type TimeHelper struct{}

func (h *TimeHelper) TimeFormat(time time.Time) string {
	weekdayFormat := h.DayFormat(time.Weekday().String())
	day := time.Day()
	monthFormat := h.MonthFormat(time.Month())
	yearFormat := time.Year()
	timeFormat := time.Format("15:04:05 WIB")

	return fmt.Sprintf("%s, %d %s %d, %s", weekdayFormat, day, monthFormat, yearFormat, timeFormat)
}

func (h *TimeHelper) DayFormat(weekday string) string {
	switch weekday {
	case "Sunday":
		return "Minggu"
	case "Monday":
		return "Senin"
	case "Tuesday":
		return "Selasa"
	case "Wednesday":
		return "Rabu"
	case "Thursday":
		return "Kamis"
	case "Friday":
		return "Jum'at"
	case "Saturday":
		return "Sabtu"
	}

	return weekday
}

func (h *TimeHelper) MonthFormat(month time.Month) string {
	switch month {
	case 1:
		return "Januari"
	case 2:
		return "Februari"
	case 3:
		return "Maret"
	case 4:
		return "April"
	case 5:
		return "Mei"
	case 6:
		return "Juni"
	case 7:
		return "Juli"
	case 8:
		return "Agustus"
	case 9:
		return "September"
	case 10:
		return "Oktober"
	case 11:
		return "November"
	case 12:
		return "Desember"
	}

	return ""
}

func NewTimeHelper() *TimeHelper {
	return &TimeHelper{}
}
