package model

type Weekday int

const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// Funktion zur Konvertierung von int zu Weekday
func IntToWeekday(day int) Weekday {
	switch day {
	case 0:
		return Monday
	case 1:
		return Tuesday
	case 2:
		return Wednesday
	case 3:
		return Thursday
	case 4:
		return Friday
	case 5:
		return Saturday
	case 6:
		return Sunday
	default:
		return Monday // Standardwert, falls die Eingabe ungültig ist
	}
}
