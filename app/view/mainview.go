package view

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/FyningTime/FyningTime/app/model"
	"github.com/FyningTime/FyningTime/app/model/db"
	"github.com/FyningTime/FyningTime/app/repo"
	"github.com/FyningTime/FyningTime/app/service"

	"github.com/charmbracelet/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppView struct {
	// Holds the main window
	window fyne.Window
	a      fyne.App

	// timetable *widget.Table represents a table widget used in the AppView struct.
	timetable *widget.Table

	// TODO add toolbar widget

	// Represents the headers for the timetable
	baseHeaders []string
	headers     []string

	// Dynamic data binding
	allOvertime binding.String

	// Actual db abstraction
	worktime  []*db.Worktime
	workday   []*db.Workday
	vacations []*db.Vacation

	cv  *CalenderView
	vpv *VacationPlannerView

	// Selected item
	selectedItem *widget.TableCellID

	// Holds the database connection
	repo *repo.SQLiteRepository
}

func (av *AppView) CreateUI(w fyne.Window, a fyne.App) *fyne.Container {
	// Assign the window to the main app struct
	av.window = w
	av.a = a

	av.allOvertime = binding.NewString()
	av.allOvertime.Set(lang.L("calculateOvertime"))

	extraColumns := len(av.baseHeaders)

	tt := widget.NewTable(
		func() (int, int) {
			// find longest day to calculate column count
			longestDay := av.getLongestWorkday()

			return len(av.workday), longestDay + extraColumns // +1 for the date row
		},
		func() fyne.CanvasObject {
			// Defines how wide the cell is
			return widget.NewLabel(lang.L("begin") + " #1")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			wd := av.workday[i.Row]
			var wtday []*db.Worktime

			for _, w := range av.worktime {
				if w.Workday.ID == wd.ID {
					wtday = append(wtday, w)
				}
			}

			// Default style
			label.TextStyle = fyne.TextStyle{Bold: false}
			label.Importance = widget.MediumImportance

			switch i.Col {
			case 0:
				label.SetText(wd.Date.Format(model.DATEFORMAT) +
					" / " + model.ShortenWeekday(wd.Date.Weekday().String()))
				label.TextStyle = fyne.TextStyle{Bold: true}
			case 1:
				label.SetText(wd.Time)
			case 2:
				label.SetText(wd.Breaktime)
			case 3:
				overtimeAsFloat, err := time.ParseDuration(wd.Overtime)
				if err != nil {
					log.Error("Overtime table field", "overtime", wd.Overtime, "error", err)
					label.SetText("")
				} else {
					if overtimeAsFloat < 0 {
						label.TextStyle = fyne.TextStyle{Bold: true}
						label.Importance = widget.HighImportance
					}
					label.SetText(wd.Overtime)
				}

			default:
				if i.Col-extraColumns < len(wtday) && i.Col > 1 {
					currentWt := wtday[i.Col-extraColumns]
					label.SetText(currentWt.Time.Format(time.TimeOnly))
				} else {
					label.SetText("")
				}
			}
		},
	)
	// Set the headers
	tt.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("")
	}
	tt.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		lbl := template.(*widget.Label)
		lbl.TextStyle = fyne.TextStyle{Bold: true}

		// Column headers (top row): Row == -1
		if id.Row == -1 && id.Col >= 0 {
			if id.Col < len(av.headers) {
				lbl.SetText(av.headers[id.Col])
			} else {
				lbl.SetText("")
			}

		}
	}
	tt.ShowHeaderRow = true
	tt.StickyColumnCount = 4

	// Date
	tt.SetColumnWidth(0, 125)
	// Time
	tt.SetColumnWidth(1, 90)
	// Break
	tt.SetColumnWidth(2, 80)
	// Overtime
	tt.SetColumnWidth(3, 95)

	tt.OnSelected = func(id widget.TableCellID) {
		av.selectedItem = &id
	}
	tt.Resize(fyne.NewSize(700, 600))

	av.timetable = tt

	btnAddTimeToolbarItem := widget.NewToolbarAction(theme.ContentAddIcon(), av.AddTimeEntry)
	btnDeleteTimeToolbarItem := widget.NewToolbarAction(theme.ContentRemoveIcon(), av.deleteButtonFunc)
	btnEditTimeToolbarItem := widget.NewToolbarAction(theme.DocumentIcon(), av.editButtonFunc)
	btnRefreshDataToolbarItem := widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
		fyne.Do(func() {
			av.RefreshData()
		})
	})

	timeToolbar := widget.NewToolbar(
		btnAddTimeToolbarItem,
		btnDeleteTimeToolbarItem,
		btnEditTimeToolbarItem,
		btnRefreshDataToolbarItem,
	)

	topBar := container.NewHBox(
		timeToolbar,
		widget.NewSeparator(),
		widget.NewLabelWithData(av.allOvertime),
		widget.NewSeparator(),
		/*widget.NewButtonWithIcon("Scroll up", theme.MoveUpIcon(), func() {
			av.timetable.ScrollToTop()
		}),*/
	)

	timerContainer := container.NewBorder(topBar, nil, nil, nil, tt)

	av.cv = CreateCalendarView(av.window, av.vacations, time.Now())
	av.vpv = CreateVacationPlannerView(av, av.repo, av.vacations)

	// Add appbar
	appTabs := container.NewAppTabs(
		container.NewTabItem(lang.L("timer"), timerContainer),
		container.NewTabItem(lang.L("calendar"), av.cv.container),
		container.NewTabItem(lang.L("vacationPlanner"), av.vpv.container),
		container.NewTabItem(lang.L("details"), widget.NewLabel(lang.L("commingSoon"))),
	)

	appContainer := container.NewBorder(nil, nil, nil, nil, appTabs)

	go av.calculateBreakLoop()
	return appContainer
}

func (av *AppView) SetBaseHeaders(headers []string) {
	av.baseHeaders = headers
}

func (av *AppView) buildHeaders() {
	// Compute the longest number of worktimes across all days
	log.Debug("Build headers for timetable")

	longest := av.getLongestWorkday()

	// Base columns
	headers := []string{}
	headers = append(headers, av.baseHeaders...)

	// Dynamic time-entry columns
	for i := range longest {
		// Option A: Label them as Begin/End pairs
		if i%2 == 0 {
			headers = append(headers, lang.L("begin")+" #"+strconv.Itoa(i/2+1))
		} else {
			headers = append(headers, lang.L("end")+" #"+strconv.Itoa(i/2+1))
		}
	}

	log.Debug("Headers built", "headers", headers)
	av.headers = headers
}

func (av *AppView) AddTimeEntry() {
	// Add a time entry to current date
	loc, _ := time.LoadLocation("Europe/Berlin")
	today, err := av.repo.GetWorkday(time.Now())
	log.Debug("Weekday", "weekday", time.Now().Weekday())
	if today == nil && err != nil {
		log.Debug("Create new workday")
		w := &db.Workday{
			Date: time.Now(),
		}
		newWd, wdErr := av.repo.AddWorkday(w)
		if wdErr != nil {
			log.Error(wdErr)
		} else {
			av.workday = append(av.workday, newWd)
			wt := &db.Worktime{
				Type:    "Begin", // As it is a new workday, it is always a begin
				Time:    time.Now().In(loc),
				Workday: *newWd,
			}

			_, wtErr := av.repo.AddWorktime(wt)
			if wtErr != nil {
				log.Error(wtErr)
			}
		}
	} else {
		log.Info("Workday already exists")
		allWt, err := av.repo.GetAllWorktime(today)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Debug("Size of worktimes", "size", len(av.worktime))
			wtType := "Begin"
			if len(allWt)%2 != 0 {
				wtType = "End"
			}

			wt := &db.Worktime{
				Type:    wtType,
				Time:    time.Now().In(loc),
				Workday: *today,
			}
			av.repo.AddWorktime(wt)
		}
	}

	// Refresh *all data*
	go av.calculateBreak(true)
}

func (av *AppView) EditSelectedTimeEntry() {
	av.editButtonFunc()
}

func (av *AppView) UnselectTableItem() {
	if av.selectedItem != nil {
		log.Debug("Unselect item", "item", av.selectedItem)
		av.timetable.Unselect(*av.selectedItem)
		av.timetable.FocusLost()
	} else {
		dialog.ShowError(errors.New(lang.L("noItemSelected")), av.window)
	}
}

func (av *AppView) DeleteSelectedTimeEntry() {
	if av.selectedItem != nil {
		_, wt, err := av.getTimeEntry(av.selectedItem)
		if err != nil {
			dialog.ShowError(err, av.window)
			return
		} else if wt == nil {
			dialog.ShowError(errors.New(lang.L("noTimeEntryFound")), av.window)
			return
		}
		log.Info("Delete selected time entry", "worktime", wt)
		defer av.deleteTimeEntry(wt)
	}
}

func (av *AppView) AddHeaders(headers []string) {
	av.headers = headers
}

func (av *AppView) CreateRepository(db *sql.DB) {
	av.repo = repo.NewSQLiteRepository(db)

	// On start we also try to migrate the database
	log.Info("Migrating database")
	err := av.repo.Migrate()
	if err != nil {
		log.Fatal(err)
	}
}

// TODO maybe limit this for a specific date range like month
func (av *AppView) RefreshData() {
	wd, wdErr := av.repo.GetAllWorkday(repo.DESC)
	if wdErr != nil {
		log.Error(wdErr)
	} else {
		av.workday = nil
		av.workday = wd
	}

	av.worktime = nil
	for i := range wd {
		wt, wtErr := av.repo.GetAllWorktime(wd[i])
		if wtErr != nil {
			log.Error(wtErr)
		} else {
			av.worktime = append(av.worktime, wt[0:]...)
		}
	}

	// Get vacations
	v, err := av.repo.GetAllVacation()
	if err != nil {
		log.Error(err)
	} else {
		av.vacations = v
		if av.vpv != nil {
			av.vpv.UpdateVacations(v)
		}
		if av.cv != nil {
			av.cv.UpdateVacations(v)
		}
	}

	av.allOvertime.Set(lang.L("calculateOvertime"))
	av.calculateOvertime()

	// Re-build headers
	av.refreshTimetable()
}

func (av *AppView) GetOvertime() string {
	totalOvertime, err := av.allOvertime.Get()
	if err != nil {
		log.Error(err)
		return ""
	}
	return totalOvertime
}

// ------------------ Private functions ------------------

func (av *AppView) deleteButtonFunc() {
	log.Info("Delete time entry", "item", av.selectedItem)
	// Check if is a day or a time entry selected
	if av.selectedItem == nil {
		dialog.ShowError(errors.New(lang.L("noItemSelected")), av.window)
		return
	} else if av.selectedItem.Col == 0 {
		dialog.ShowConfirm(lang.L("deleteEntry"), lang.L("areYouSureDelete"), func(b bool) {
			if b {
				// Delete the whole day
				wd := av.workday[av.selectedItem.Row]
				log.Info("Delete workday", "workday", wd)
				// Deletion here
				rows, err := av.repo.DeleteWorkday(wd)
				if rows != 0 || err == nil {
					// Refresh *all data*
					av.refreshAll()
				} else {
					dialog.ShowError(errors.New(lang.L("couldNotDeleteDataset")), av.window)
				}
			}
		}, av.window)
	} else if av.selectedItem.Col == 1 {
		dialog.ShowError(errors.New(lang.L("noItemSelected")), av.window)
	} else {
		// Assure that there is an item selected
		_, wt, err := av.getTimeEntry(av.selectedItem)

		// Delete only if there is an existing item and no error
		// Else show an error dialog
		if wt != nil && err == nil {
			dialog.ShowConfirm(lang.L("deleteEntry"), lang.L("areYouSureDelete"), func(b bool) {
				if b {
					defer av.deleteTimeEntry(wt)
				} else {
					log.Debug("Delete time entry", "canceled", b)
				}
			}, av.window)
		} else {
			dialog.ShowError(errors.New(lang.L("noItemSelected")), av.window)
		}
	}

	go av.calculateBreak(true)
}

func (av *AppView) getTimeEntry(item *widget.TableCellID) (*db.Workday, *db.Worktime, error) {
	log.Info("Get time entry", "item", item)
	if item == nil {
		return nil, nil, errors.New(lang.L("noItemSelected"))
	}

	// Get the affacted workday by row
	wd := av.workday[av.selectedItem.Row]
	log.Debug("Get time entry", "workday", wd)

	// Get all worktimes for this workday
	var wtList []*db.Worktime
	for _, w := range av.worktime {
		if w.Workday.ID == wd.ID {
			log.Debug("Found worktime", "worktime", w)
			wtList = append(wtList, w)
		}
	}

	extraColumns := len(av.baseHeaders)

	// Get the worktime by column (-2 because of the date column)
	if av.selectedItem.Col > (extraColumns-1) && av.selectedItem != nil &&
		av.selectedItem.Col-extraColumns < len(wtList) {
		return wd, wtList[av.selectedItem.Col-extraColumns], nil
	} else {
		return nil, nil, errors.New(lang.L("noWorktimeFound"))
	}
}

// TODO Important! simplify this function
func (av *AppView) editButtonFunc() {
	loc, _ := time.LoadLocation("Europe/Berlin")

	// Declare if we have to add or update a time entry
	isAdd := false

	log.Info("Edit time entry", "item", av.selectedItem)
	// Check if is a day or a time entry selected
	// Assure that there is an item selected
	if av.selectedItem == nil || av.selectedItem.Col == 0 || av.selectedItem.Col == 1 {
		dialog.ShowError(errors.New(lang.L("noWorktimeSelected")), av.window)
		return
	}

	wd, wt, err := av.getTimeEntry(av.selectedItem)
	if err != nil && wt == nil {
		// If no worktime is found, we try to use the previous worktime. Maybe it was forgot to end the worktime
		av.selectedItem.Col--

		log.Debug("Get fallback time entry", "item", av.selectedItem)
		wd, wt, err = av.getTimeEntry(av.selectedItem)
		if err != nil {
			dialog.ShowError(err, av.window)
			return
		}
		isAdd = true
		if wt.Type == "End" {
			wt.Type = "Begin"
		} else {
			wt.Type = "End"
		}
	} else if wt.Workday.Date.Before(time.Now().In(loc)) {
		isAdd = false
	}

	log.Debug("Is add worktime", "isAdd", isAdd)
	if wt != nil && err == nil {
		timeEntry := widget.NewEntry()
		timeEntry.SetText(time.Now().In(loc).Format(time.TimeOnly))

		typeEntry := widget.NewSelectEntry([]string{lang.L("begin"), lang.L("end")})
		var textType string
		if wt.Type == "Begin" {
			textType = lang.L("begin")
		} else {
			textType = lang.L("end")
		}

		typeEntry.SetText(textType)
		typeEntry.Disable()

		form := []*widget.FormItem{
			{Text: lang.L("timeWithFormat"),
				Widget:   timeEntry,
				HintText: lang.L("hintTextOldEntry") + ": " + wt.Time.Format(time.TimeOnly)},
			{Text: lang.L("type"), Widget: typeEntry},
		}
		// Edit only if there is an existing item and no error
		log.Debug("Edit time entry", "worktime", wt)
		dia := dialog.NewForm(lang.L("editTimeEntry"), lang.L("edit"), lang.L("cancel"), form, func(b bool) {
			if b {
				log.Debug("Edit time entry", "confirmed", b, "worktime", wt)
				// Parse the time
				timeEntry := form[0].Widget.(*widget.Entry).Text

				if timeEntry == "" {
					dialog.ShowError(errors.New(lang.L("noTimeEntry")), av.window)
					return
				}

				// Add the current date to the updated time entry
				tempTime := wd.Date.In(loc).Format(time.DateOnly) + " " + timeEntry
				nt, err := time.ParseInLocation(time.DateTime, tempTime, loc)
				if err != nil {
					dialog.ShowError(err, av.window)
					return
				}
				log.Info("New time", "time", nt)

				prevCol := widget.TableCellID{
					Row: av.selectedItem.Row,
					Col: av.selectedItem.Col - 1,
				}
				if prevCol.Col >= 1 { // If there is a previous entry
					wd, prevTime, err := av.getTimeEntry(&prevCol)
					log.Debug("Previous time", "time", "workday", prevTime.Time, wd)
					if err != nil {
						dialog.ShowError(err, av.window)
						return
					}
				}

				// The edit shouldn't be in future, it would be faking and does not make sense
				ct := time.Now().In(loc)
				log.Debug("Compare new-time with current-time", "new-time", nt, "current-time", ct)
				if ct.Before(nt) {
					dialog.ShowError(errors.New(lang.L("timeIsInFuture")), av.window)
					return
				}

				// Update the time entry
				wt.Time = nt
				timeType := form[1].Widget.(*widget.SelectEntry).Text

				if timeType == lang.L("begin") {
					wt.Type = "Begin"
				} else {
					wt.Type = "End"
				}

				if isAdd {
					// If we are adding a time entry in the past, we have to add a new workday
					_, errAdd := av.repo.AddWorktime(wt)
					if errAdd != nil {
						dialog.ShowError(err, av.window)
					}
				} else {
					_, errUpdate := av.repo.UpdateWorktime(wt)
					if errUpdate != nil {
						dialog.ShowError(err, av.window)
					}
				}
			}
			// Refresh *all data*
			go av.calculateBreak(true)
		}, av.window)

		dia.Resize(fyne.NewSize(400, 200))
		dia.Show()
	} else {
		dialog.ShowError(errors.New(lang.L("noItemSelected")), av.window)
	}
}

func (av *AppView) calculateBreakLoop() {
	for {
		av.calculateBreak(false)
	}
}

/*
Calculates the breaktime for a workday
and refreshes the data
*/
func (av *AppView) calculateBreak(skipWait ...bool) {
	// Run this all time in the background
	settings := service.ReadProperties(av.a)

	if len(skipWait) > 0 && !skipWait[0] {
		// Wait until the time to sleep
		tts := time.Duration(settings.RefreshTimeUi) * time.Second
		log.Debug("Wait to calculate breaktime", "time-to-sleep", tts)
		time.Sleep(tts)
	}
	log.Info("Calculate breaktime")

	wd, errWd := av.repo.GetAllWorkday(repo.DESC)
	if errWd != nil {
		log.Error(errWd)
	}

	for _, w := range wd {
		log.Debug("Workday", "workday", w)
		wts, errWt := av.repo.GetAllWorktime(w)
		if errWt != nil {
			log.Error(errWt)
		}

		/*
			breaktime := 0min when time <= 6:00
			breaktime := 30min when time > 6:00
			breaktime := 45min when time > 9:00
		*/
		var breaktime time.Duration = 0 * time.Minute
		var worktime time.Duration = 0 * time.Minute
		var workTimeDiff time.Duration = 0 * time.Minute

		wtsLength := len(wts)
		tempWtsLength := wtsLength
		tempWts := wts

		if tempWtsLength%2 == 0 && tempWtsLength >= 4 {
			for tempWtsLength > 2 && tempWtsLength%2 == 0 {
				// Calculate the difference between begin and the previous end time
				lastBegin := tempWts[tempWtsLength-2]
				lastEnd := tempWts[tempWtsLength-3]
				workTimeDiff = workTimeDiff + lastBegin.Time.Sub(lastEnd.Time)
				log.Debug("Worktime difference for even entries", "worktime-diff", workTimeDiff)
				tempWtsLength -= 2
				tempWts = tempWts[:tempWtsLength]
			}
		}

		for c, wt := range wts {
			if c+1 < wtsLength {
				if wt.Type == "Begin" {
					end := wts[c+1]
					// Entfernen der Nanosekunden von der Zeit
					startTime := wt.Time.Truncate(time.Second)
					endTime := end.Time.Truncate(time.Second)
					worktime += endTime.Sub(startTime)
				}
			}
		}

		log.Debug("All day worktime", "worktime", worktime)

		if workTimeDiff == 0 {
			switch {
			case worktime < 6*time.Hour:
				breaktime = 0 * time.Minute
			case worktime >= 6*time.Hour && worktime < 9*time.Hour:
				breaktime = 30 * time.Minute
			case worktime >= 9*time.Hour:
				breaktime = 45 * time.Minute
			default: // This should never happen
				breaktime = 30 * time.Minute
			}
		} else {
			breaktime = workTimeDiff
		}

		// FIXME currently it's saved as string in the database as workaround
		// Update the workday with the calculated breaktime
		w.Breaktime = breaktime.Truncate(time.Minute).String()
		if strings.Contains(w.Breaktime, "m0s") {
			w.Breaktime = w.Breaktime[:len(w.Breaktime)-2]
		}
		if workTimeDiff == 0 {
			w.Time = (worktime - breaktime).String()
		} else {
			w.Time = worktime.String()
		}
		log.Debug("Update workday",
			"worktime", worktime, "breaktime", breaktime)

		_, errUpdate := av.repo.UpdateWorkday(w)
		if errUpdate != nil {
			log.Error(errUpdate)
		}
	}

	// Refresh *all data*
	av.refreshAll()
}

/*
Deletes a time entry from the database
*/
func (av *AppView) deleteTimeEntry(wt *db.Worktime) {
	// Deletion here
	rows, err := av.repo.DeleteWorktime(wt)
	if rows != 0 || err != nil {
		// Refresh *all data*
		av.refreshAll()
	}
}

func (av *AppView) refreshTimetable() {
	av.buildHeaders()
	av.timetable.Refresh()
}

/*
Combines the refresh of the data and the timetable
*/
func (av *AppView) refreshAll() {
	v, err := av.repo.GetAllVacation()
	if err != nil {
		log.Error(err)
	} else {
		av.vacations = v
	}
	fyne.Do(func() {
		// Workaround
		av.RefreshData()
		av.calculateOvertime()
		av.RefreshData()
	})
}

func (av *AppView) getLongestWorkday() int {
	var longestDay int = 0
	var tempDay int = 0

	// Find the longest workday
	for _, wd := range av.workday {
		tempDay = 0
		for _, w := range av.worktime {
			if w.Workday.ID == wd.ID {
				tempDay++
			}
		}

		if tempDay > longestDay {
			longestDay = tempDay
		}
	}
	return longestDay
}

// Slice is a workaround as there is no optional parameters in Go
func (av *AppView) calculateOvertime(previousOvertime ...time.Duration) {
	log.Debug("Calculate overtime")

	settings := service.ReadProperties(av.a)

	wd, err := av.repo.GetAllWorkday(repo.ASC)
	if err != nil {
		log.Error(err)
		dialog.ShowError(err, av.window)
	} else {
		var totalOvertime time.Duration = 0 * time.Hour
		var previousOvertimeTransfered time.Duration = 0 * time.Hour
		var workHoursPerDayDuration time.Duration = getHoursPerDay(settings.WeekHours)
		var importTotalOvertime time.Duration = time.Duration(settings.ImportOvertime * float64(time.Hour))

		// Fake it till you make it
		if len(previousOvertime) > 0 {
			previousOvertimeTransfered = previousOvertime[0]
		}
		log.Debug("Work hour per day", "duration", workHoursPerDayDuration.String())

		for _, w := range wd {
			currentWorktimeDb, err := time.ParseDuration(w.Time)

			midSumOvertime := (currentWorktimeDb - workHoursPerDayDuration)
			if err != nil {
				log.Error(err)
				dialog.ShowError(err, av.window)
			} else {
				w.Overtime = midSumOvertime.String()
				// Defer DB update to batch after loop
				totalOvertime += midSumOvertime
			}
		}
		// Batch update all overtimes in one DB call
		if err := av.repo.UpdateOvertimesBatch(wd); err != nil {
			log.Error("Batch update of overtimes failed", "error", err)
			dialog.ShowError(err, av.window)
		}

		// Add previous overtime transfered from last calculation
		totalOvertime += previousOvertimeTransfered

		log.Debug("Previous overtime transfered", "overtime", previousOvertimeTransfered)

		// Add imported overtime from settings
		if importTotalOvertime > 0 {
			totalOvertime += importTotalOvertime
			log.Info("Import overtime", "overtime", importTotalOvertime)
		}

		if av.allOvertime == nil {
			av.allOvertime = binding.NewString()
		}

		av.allOvertime.Set(lang.L("totalOvertime") + ": " + (totalOvertime + previousOvertimeTransfered).String())
		log.Info("Total overtime", "overtime", totalOvertime+previousOvertimeTransfered)
	}
}

// TODO Move logic from editButtonFunc here
func (av *AppView) editSingleTimeEntry() {

}

func (av *AppView) editDateEntry() {

}

// Calculate hours per day, assuming 5 working days per week
func getHoursPerDay(weekHours int) time.Duration {
	weekHoursPerDay := float64(weekHours) / float64(5)
	return time.Duration(weekHoursPerDay * float64(time.Hour))
}
