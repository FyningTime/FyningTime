package service

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/FyningTime/FyningTime/app/model"
)

func ReadSettings() (*model.Settings, error) {
	settingsFilePath, err := GetFyningTimePath()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(settingsFilePath)
	if err != nil && file == nil {
		log.Error(err)
		// Create settings file if it doesn't exist
		log.Infof("Creating settings file, %s ", settingsFilePath)

		dbPath, err := GetFyningTimePath(model.DBFILE)
		if err != nil {
			log.Fatal(err)
		}

		settings := model.NewSettings(
			settingsFilePath, dbPath,
		)

		createError := WriteSettings(settings)
		if createError != nil {
			log.Error(createError)
			return nil, createError
		}
		file, _ = os.Open(settingsFilePath)
	}
	defer file.Close()

	settings := model.Settings{}
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	return validateSettings(&settings)
}

func WriteSettings(s *model.Settings) error {
	log.Debugf("Writing settings: %+v\n", s)
	validateSettings(s)
	settingsFilePath, err := GetFyningTimePath()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(settingsFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		newFile, err := os.Create(settingsFilePath)
		if err != nil {
			log.Fatal(err)
			return err
		}
		file = newFile
	}
	//defer file.Close()

	byteValue, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = file.Write(byteValue)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	return nil
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
