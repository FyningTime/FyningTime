package view

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/FyningTime/FyningTime/app/model"
	"github.com/FyningTime/FyningTime/app/service"
)

func GetSettingsView(w fyne.Window) *dialog.FormDialog {
	settings, err := service.ReadSettings()
	if err != nil {
		dialog.ShowError(err, w)
		return nil
	}

	firstDayOfWeekEntry := widget.NewSelectEntry(
		[]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"})
	weekday := settings.FirstDayOfWeek
	firstDayOfWeekEntry.SetText(model.WeekdayToString(weekday))

	maxVacations := widget.NewEntry()
	maxVacations.SetText(strconv.Itoa(settings.MaxVacationDays))

	refreshTimeUi := widget.NewEntry()
	refreshTimeUi.SetText(strconv.Itoa(settings.RefreshTimeUi))

	form := []*widget.FormItem{
		{Text: "Saved path", Widget: widget.NewLabel(settings.SavedPath)},
		{Text: "DB path", Widget: widget.NewLabel(settings.SavedDbPath)},
		{Text: "Refresh times (seconds)", Widget: refreshTimeUi},
		{Text: "First day of week", Widget: firstDayOfWeekEntry},
		{Text: "Max vacations", Widget: maxVacations},
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
			settings.RefreshTimeUi, err = strconv.Atoi(refreshTimeUi.Text)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			service.WriteSettings(settings)
		}
	}, w)
	return dia
}
