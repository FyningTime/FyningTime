package service

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"github.com/FyningTime/FyningTime/app/model"
)

const (
	SETTINGSFILE string = "settings.json"
)

func ReadSettings() (*model.Settings, error) {
	settingsFilePath, err := getSettingsFilePath()
	if err != nil {
		log.Fatal(err)
	}

	settings := model.Settings{}

	// Join the home directory with the settings file name
	settings.SavedPath = settingsFilePath

	file, err := os.Open(settingsFilePath)
	if err != nil && file == nil {
		log.Error(err)
		// Create settings file if it doesn't exist
		log.Info("Creating settings file: " + settingsFilePath)
		createError := WriteSettings(&settings)
		if createError != nil {
			log.Error(createError)
			return nil, createError
		}
		file, _ = os.Open(settingsFilePath)
	}
	defer file.Close()

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

	createError := WriteSettings(&settings)
	if createError != nil {
		log.Error(createError)
		return nil, createError
	}

	return &settings, nil
}

func WriteSettings(s *model.Settings) error {
	settingsFilePath, err := getSettingsFilePath()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(settingsFilePath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	byteValue, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = file.Write(byteValue)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func getSettingsFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	settingsFilePath := filepath.Join(homeDir, ".fyningtime", SETTINGSFILE)
	path := filepath.Dir(settingsFilePath)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return settingsFilePath, nil
}
