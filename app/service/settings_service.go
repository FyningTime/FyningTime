package service

import (
	"errors"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"github.com/charmbracelet/log"

	"github.com/FyningTime/FyningTime/app/model"
)

type SettingsProperty struct {
	property string
	value    any
}

const (
	// Settings property names
	weekHoursProperty       = "weekHours"
	firstDayOfWeekProperty  = "firstDayOfWeek"
	maxVacationDaysProperty = "maxVacationDays"
	importOvertimeProperty  = "importOvertime"
	refreshTimeUiProperty   = "refreshTimeUi"
	themeVariantProperty    = "themeVariant"

	// Default settings values
	weekHoursDefault       = 40
	firstDayOfWeekDefault  = model.Monday // Monday
	maxVacationDaysDefault = 30
	importOvertimeDefault  = 0
	refreshTimeUiDefault   = 300 // in seconds
	themeVariantDefault    = 0   // 0=light/auto, 1=dark
)

func ReadProperties(a fyne.App) *model.Settings {
	dbPath, err := GetFyningTimePath(model.DBFILE)
	if err != nil {
		log.Fatal(err)
	}

	settings := model.NewSettings(
		"", dbPath,
	)
	settings.WeekHours = a.Preferences().IntWithFallback(weekHoursProperty, weekHoursDefault)
	settings.FirstDayOfWeek = model.Weekday(a.Preferences().StringWithFallback(firstDayOfWeekProperty, model.WeekdayToString(firstDayOfWeekDefault)))
	settings.MaxVacationDays = a.Preferences().IntWithFallback(maxVacationDaysProperty, maxVacationDaysDefault)
	settings.ImportOvertime = a.Preferences().FloatWithFallback(importOvertimeProperty, importOvertimeDefault)
	settings.RefreshTimeUi = a.Preferences().IntWithFallback(refreshTimeUiProperty, refreshTimeUiDefault)
	settings.ThemeVariant = a.Preferences().IntWithFallback(themeVariantProperty, themeVariantDefault)
	return settings
}

func WriteProperties(a fyne.App, s *model.Settings) {
	a.Preferences().SetInt("weekHours", s.WeekHours)
	a.Preferences().SetString("firstDayOfWeek", model.WeekdayToString(s.FirstDayOfWeek))
	a.Preferences().SetInt("maxVacationDays", s.MaxVacationDays)
	a.Preferences().SetFloat("importOvertime", s.ImportOvertime)
	a.Preferences().SetInt("refreshTimeUi", s.RefreshTimeUi)
	a.Preferences().SetInt("themeVariant", s.ThemeVariant)
}

/**
 * Get the path to the fyning file
 * Default it returns the path to the settings file
 */
func GetFyningTimePath(file ...string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fileName := model.SETTINGSFILE
	if len(file) > 0 && file[0] != "" {
		fileName = file[0]
	}

	fyningPath := filepath.Join(homeDir, model.FYNINGTIMEDIR)
	err = os.MkdirAll(fyningPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fyningFile := filepath.Join(fyningPath, fileName)
	return fyningFile, nil
}

/**
 * Get the path to the fyning file
 * Default it returns the path to the settings file
 */
func validateSettings(s *model.Settings) (*model.Settings, error) {
	var err error

	// Paths should be automatically set by application
	// it's just for information where the files are stored
	if s.SavedDbPath == "" {
		log.Error("DB path is not set correctly")
		err = errors.New("DB path is not set correctly")
	}
	if s.SavedPath == "" {
		log.Error("Saved path is not set correctly")
		err = errors.New("saved path is not set correctly")
	}

	// UI specific configuration
	// Refresh rate to determin work breaks should be at least 15 seconds
	if s.RefreshTimeUi < 15 { // 15 seconds is the minimum
		log.Warn("Refresh time is not set correctly. Set to 15 seconds.")
		s.RefreshTimeUi = 15
	}

	// Business logic specific configuration

	if s.WeekHours < 0 { // 0 is the minimum
		log.Warn("Week hours is not set correctly. Set to 40 hours.")
		s.WeekHours = 40
	}

	// In Europe the first work day of the week is Monday
	if s.FirstDayOfWeek < model.Monday || s.FirstDayOfWeek > model.Sunday {
		log.Warn("First day of week is not set correctly. Set to Monday.")
		s.FirstDayOfWeek = model.Monday
	}
	// Everyone should have vacations
	if s.MaxVacationDays < 0 { // 0 is the minimum
		log.Warn("Max vacation days is not set correctly. Set to 30 days.")
		s.MaxVacationDays = 30
	}

	if s.ImportOvertime < 0 { // 0 is the minimum
		log.Warn("Import overtime is not set correctly. Set to 0 hours.")
		s.ImportOvertime = 0
	}

	return s, err
}
