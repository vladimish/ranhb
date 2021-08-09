package date

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
