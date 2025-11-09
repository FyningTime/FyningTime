package view

import (
	"errors"

	"strconv"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/dialog"

	"fyne.io/fyne/v2/widget"

	"github.com/FyningTime/FyningTime/app/model"

	"github.com/FyningTime/FyningTime/app/service"
	apptheme "github.com/FyningTime/FyningTime/app/theme"
)

func GetSettingsView(w fyne.Window, a fyne.App) *dialog.FormDialog {
	settings := service.ReadProperties(a)

	firstDayOfWeekEntry := widget.NewSelectEntry(
		[]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"})
	weekday := settings.FirstDayOfWeek
	firstDayOfWeekEntry.SetText(model.WeekdayToString(weekday))

	weekHours := widget.NewEntry()
	weekHours.SetText(strconv.Itoa(settings.WeekHours))

	importTotalOvertime := widget.NewEntry()
	importTotalOvertime.SetText(strconv.FormatFloat(settings.ImportOvertime, 'f', -1, 64))

	maxVacations := widget.NewEntry()
	maxVacations.SetText(strconv.Itoa(settings.MaxVacationDays))

	refreshTimeUi := widget.NewEntry()

	refreshTimeUi.SetText(strconv.Itoa(settings.RefreshTimeUi))

	themeOptions := []string{"Auto", "Light", "Dark"}
	themeSelection := widget.NewRadioGroup(themeOptions, nil)
	switch settings.ThemeVariant {
	case 1:
		themeSelection.SetSelected("Dark")
	case 2:
		themeSelection.SetSelected("Light")
	default:
		themeSelection.SetSelected("Auto")
	}

	form := []*widget.FormItem{
		{Text: "Saved path", Widget: widget.NewLabel(settings.SavedPath)},
		{Text: "DB path", Widget: widget.NewLabel(settings.SavedDbPath)},
		{Text: "Refresh times (seconds)", Widget: refreshTimeUi},
		{Text: "First day of week", Widget: firstDayOfWeekEntry},
		{Text: "Week hours", Widget: weekHours},
		{Text: "Max vacations", Widget: maxVacations},
		{Text: "Import total overtime", Widget: importTotalOvertime},
		{Text: "Theme", Widget: themeSelection},
	}
	dia := dialog.NewForm("Settings", "Save", "Cancel", form, func(ok bool) {
		if ok {
			// Save settings
			settings.FirstDayOfWeek = model.StringToWeekday(firstDayOfWeekEntry.Text)

			intMaxVacations, err := strconv.Atoi(maxVacations.Text)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			settings.MaxVacationDays = intMaxVacations

			intWeekHours, err := strconv.Atoi(weekHours.Text)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if intWeekHours < 1 || intWeekHours > 50 {
				dialog.ShowError(errors.New("Week hours must be between 1 and 50"), w)
				return
			}
			settings.WeekHours = intWeekHours

			intImportOvertime, err := strconv.ParseFloat(importTotalOvertime.Text, 64)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if intImportOvertime < 0 {
				dialog.ShowError(errors.New("Import total overtime cannot be negative"), w)
				return
			}
			settings.ImportOvertime = intImportOvertime

			settings.RefreshTimeUi, err = strconv.Atoi(refreshTimeUi.Text)

			if err != nil {

				dialog.ShowError(err, w)

				return

			}

			switch themeSelection.Selected {
			case "Dark":
				settings.ThemeVariant = 1
				a.Settings().SetTheme(apptheme.NewPastelleTheme())
			case "Light":
				settings.ThemeVariant = 2
				a.Settings().SetTheme(apptheme.NewPastelleLight())
			default:
				settings.ThemeVariant = 0
				a.Settings().SetTheme(apptheme.NewPastelleTheme())
			}
			service.WriteProperties(a, settings)

		} else {
			// Canceled
		}
	}, w)
	return dia
}
