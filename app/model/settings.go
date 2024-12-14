package model

type Settings struct {
	SavedPath       string  `json:"saved_path"`
	SavedDbPath     string  `json:"saved_db_path"`
	FirstDayOfWeek  Weekday `json:"first_day_of_week"`
	MaxVacationDays int     `json:"max_vacation_days"`
	RefreshTimeUi   int     `json:"refresh_time_ui"`
	ThemeVariant    int     `json:"theme_variant"`
}

func NewSettings(savedPath string, savedDbPath string) *Settings {
	return &Settings{
		SavedPath:       savedPath,
		SavedDbPath:     savedDbPath,
		FirstDayOfWeek:  Monday,
		MaxVacationDays: 30,
		// Refresh time in seconds
		RefreshTimeUi: 30,
	}
}
