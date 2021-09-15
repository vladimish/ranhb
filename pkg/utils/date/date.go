package date

import "time"

func IntToWeekday(i int) string {
	switch i {
	case 7:
		return "Воскресенье"
	case 1:
		return "Понедельник"
	case 2:
		return "Вторник"
	case 3:
		return "Среда"
	case 4:
		return "Четверг"
	case 5:
		return "Пятница"
	case 6:
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
