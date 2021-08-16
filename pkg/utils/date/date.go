package date

import "time"

func IntToWeekday(i int) string {
	switch i {
	case 1:
		return "Воскресенье"
	case 2:
		return "Понедельник"
	case 3:
		return "Вторник"
	case 4:
		return "Среда"
	case 5:
		return "Четверг"
	case 6:
		return "Пятница"
	case 7:
		return "Суббота"
	default:
		return ""
	}
}

// GetWeekInterval returns two times
// which corresponds for monday and sunday
// of a d week.
func GetWeekInterval(d time.Time) (startTime time.Time, endTime time.Time) {
	for d.Weekday() != time.Monday {
		d = d.AddDate(0, 0, -1)
	}
	startTime = d
	endTime = startTime.Add(6 * 24 * time.Hour)

	return startTime, endTime
}