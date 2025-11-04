package model

type Weekday string

const (
	Monday    Weekday = "Monday"
	Tuesday   Weekday = "Tuesday"
	Wednesday Weekday = "Wednesday"
	Thursday  Weekday = "Thursday"
	Friday    Weekday = "Friday"
	Saturday  Weekday = "Saturday"
	Sunday    Weekday = "Sunday"
)

func StringToWeekday(day string) Weekday {
	switch day {
	case "Monday":
		return Monday
	case "Tuesday":
		return Tuesday
	case "Wednesday":
		return Wednesday
	case "Thursday":
		return Thursday
	case "Friday":
		return Friday
	case "Saturday":
		return Saturday
	case "Sunday":
		return Sunday
	default:
		return Monday // Standardwert, falls die Eingabe ungültig ist
	}
}

func WeekdayToString(day Weekday) string {
	switch day {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	default:
		return "Monday" // Standardwert, falls die Eingabe ungültig ist
	}
}

func ShortenWeekday(day string) string {
	switch day {
	case "Monday":
		return "Mon"
	case "Tuesday":
		return "Tue"
	case "Wednesday":
		return "Wed"
	case "Thursday":
		return "Thu"
	case "Friday":
		return "Fri"
	case "Saturday":
		return "Sat"
	case "Sunday":
		return "Sun"
	default:
		return "Mon" // Standardwert, falls die Eingabe ungültig ist
	}
}
