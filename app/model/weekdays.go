package model

import "fyne.io/fyne/v2/lang"

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
	case lang.L("monday"):
		return Monday
	case lang.L("tuesday"):
		return Tuesday
	case lang.L("wednesday"):
		return Wednesday
	case lang.L("thursday"):
		return Thursday
	case lang.L("friday"):
		return Friday
	case lang.L("saturday"):
		return Saturday
	case lang.L("sunday"):
		return Sunday
	default:
		return Monday // Standardwert, falls die Eingabe ungültig ist
	}
}

func WeekdayToString(day Weekday) string {
	switch day {
	case Monday:
		return lang.L("monday")
	case Tuesday:
		return lang.L("tuesday")
	case Wednesday:
		return lang.L("wednesday")
	case Thursday:
		return lang.L("thursday")
	case Friday:
		return lang.L("friday")
	case Saturday:
		return lang.L("saturday")
	case Sunday:
		return lang.L("sunday")
	default:
		return lang.L("monday") // Standardwert, falls die Eingabe ungültig ist
	}
}

func ShortenWeekday(day string) string {
	switch day {
	case "Monday":
		return lang.L("mondayShort")
	case "Tuesday":
		return lang.L("tuesdayShort")
	case "Wednesday":
		return lang.L("wednesdayShort")
	case "Thursday":
		return lang.L("thursdayShort")
	case "Friday":
		return lang.L("fridayShort")
	case "Saturday":
		return lang.L("saturdayShort")
	case "Sunday":
		return lang.L("sundayShort")
	default:
		return lang.L("mondayShort") // Standardwert, falls die Eingabe ungültig ist
	}
}
