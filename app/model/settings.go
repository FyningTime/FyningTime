package model

type Settings struct {
	FirstDayOfWeek  Weekday `json:"first_day_of_week"`
	SavedPath       string  `json:"saved_path"`
	MaxVacationDays int     `json:"max_vacation_days"`
}
