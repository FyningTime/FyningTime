package model

type Settings struct {
	SavedPath   string `json:"saved_path"`
	SavedDbPath string `json:"saved_db_path"`
	// UI specific configuration
	RefreshTimeUi int `json:"refresh_time_ui"`
	ThemeVariant  int `json:"theme_variant"`

	// Business logic specific configuration
	FirstDayOfWeek Weekday `json:"first_day_of_week"`
	// How many hours a week should be worked
	// Necessary for overtime working
	WeekHours       int `json:"week_hours"`
	MaxVacationDays int `json:"max_vacation_days"`
}

func NewSettings(savedPath string, savedDbPath string) *Settings {
	return &Settings{
		SavedPath:   savedPath,
		SavedDbPath: savedDbPath,

		// UI specific configuration
		// Refresh time in seconds
		RefreshTimeUi: 30,
		ThemeVariant:  0,

		// Business logic specific configuration
		FirstDayOfWeek:  Monday,
		MaxVacationDays: 30,
		WeekHours:       40,
	}
}
