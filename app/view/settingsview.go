package view

import (
	"errors"

	"strconv"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"

	"fyne.io/fyne/v2/widget"

	"github.com/FyningTime/FyningTime/app/model"

	"github.com/FyningTime/FyningTime/app/service"
	apptheme "github.com/FyningTime/FyningTime/app/theme"
)

func GetSettingsView(w fyne.Window, a fyne.App) *dialog.FormDialog {
	settings := service.ReadProperties(a)

	firstDayOfWeekEntry := widget.NewSelectEntry(
		[]string{
			lang.L("monday"),
			lang.L("tuesday"),
			lang.L("wednesday"),
			lang.L("thursday"),
			lang.L("friday"),
			lang.L("saturday"),
			lang.L("sunday"),
		})
	weekday := settings.FirstDayOfWeek
	firstDayOfWeekEntry.SetText(model.WeekdayToString(weekday))

	weekHours := widget.NewEntry()
	weekHours.SetText(strconv.Itoa(settings.WeekHours))

	lockImportOvertime := widget.NewCheck(lang.L("lockImportOvertime"), nil)
	lockImportOvertime.SetChecked(settings.LockImportOvertime)

	importTotalOvertime := widget.NewEntry()
	importTotalOvertime.SetText(strconv.FormatFloat(settings.ImportOvertime, 'f', -1, 64))
	if settings.LockImportOvertime {
		importTotalOvertime.Disable()
	}

	maxVacations := widget.NewEntry()
	maxVacations.SetText(strconv.Itoa(settings.MaxVacationDays))

	refreshTimeUi := widget.NewEntry()

	refreshTimeUi.SetText(strconv.Itoa(settings.RefreshTimeUi))

	themeOptions := []string{lang.L("auto"), lang.L("light"), lang.L("dark")}
	themeSelection := widget.NewRadioGroup(themeOptions, nil)
	switch settings.ThemeVariant {
	case 1:
		themeSelection.SetSelected(lang.L("dark"))
	case 2:
		themeSelection.SetSelected(lang.L("light"))
	default:
		themeSelection.SetSelected(lang.L("auto"))
	}

	form := []*widget.FormItem{
		{Text: lang.L("dbPath"), Widget: widget.NewLabel(settings.SavedDbPath)},
		{Text: lang.L("refreshTimesInSeconds"), Widget: refreshTimeUi},
		{Text: lang.L("firstDayOfWeek"), Widget: firstDayOfWeekEntry},
		{Text: lang.L("weekHours"), Widget: weekHours},
		{Text: lang.L("maxVacations"), Widget: maxVacations},
		{Text: lang.L("importTotalOvertime"), Widget: importTotalOvertime},
		{Text: lang.L("theme"), Widget: themeSelection},
		{Text: lang.L("lockImportOvertime"), Widget: lockImportOvertime},
	}
	dia := dialog.NewForm(lang.L("settings"), lang.L("save"), lang.L("cancel"), form, func(ok bool) {
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
			settings.LockImportOvertime = lockImportOvertime.Checked

			settings.RefreshTimeUi, err = strconv.Atoi(refreshTimeUi.Text)

			if err != nil {
				dialog.ShowError(err, w)
				return

			}

			switch themeSelection.Selected {
			case lang.L("dark"):
				settings.ThemeVariant = 1
				a.Settings().SetTheme(apptheme.NewPastelleDark())
			case lang.L("light"):
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
