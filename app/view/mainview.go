package view

import (
	"database/sql"
	"errors"
	"main/app/model/db"
	"main/app/repo"
	"strconv"
	"time"

	"github.com/charmbracelet/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppView struct {
	// Holds the main window
	window fyne.Window

	// timetable *widget.Table represents a table widget used in the AppView struct.
	timetable *widget.Table

	// TODO add toolbar widget

	// Represents the headers for the timetable
	headers []string

	// Actual db abstraction
	worktime []*db.Worktime
	workday  []*db.Workday

	// Selected item
	selectedItem *widget.TableCellID

	// Holds the database connection
	repo *repo.SQLiteRepository
}

func (av *AppView) CreateUI(w fyne.Window) *container.Split {
	// Assign the window to the main app struct
	av.window = w

	tt := widget.NewTable(
		func() (int, int) {
			// find longest day to calculate column count
			var longestDay int
			var tempDay int

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

			return len(av.workday), longestDay + 1 // +1 for the date row
		},
		func() fyne.CanvasObject {
			// Defines how wide the cell is
			return widget.NewLabel("15:03:21 (Begin)")
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

			if i.Col == 0 {
				label.SetText(wd.Date.Format(time.DateOnly) +
					" id: " + strconv.FormatInt(wd.ID, 10))
			} else {
				if i.Col-1 < len(wtday) {
					currentWt := wtday[i.Col-1]
					workType := " (" + currentWt.Type + ") "
					label.SetText(currentWt.Time.Format(time.TimeOnly) + workType)
				} else {
					label.SetText("")
				}
			}
		},
	)
	// Set the headers
	tt.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("Date")
	}
	tt.UpdateHeader = func(i widget.TableCellID, o fyne.CanvasObject) {
	}
	tt.OnSelected = func(id widget.TableCellID) {
		av.selectedItem = &id
	}
	tt.Resize(fyne.NewSize(1000, 600))
	av.timetable = tt

	btnAddTimeToolbarItem := widget.NewToolbarAction(theme.ContentAddIcon(), av.AddTimeEntry)
	btnDeleteTimeToolbarItem := widget.NewToolbarAction(theme.ContentRemoveIcon(), av.deleteButtonFunc)
	btnEditTimeToolbarItem := widget.NewToolbarAction(theme.ContentCutIcon(), av.editButtonFunc)

	timeToolbar := widget.NewToolbar(btnAddTimeToolbarItem,
		btnDeleteTimeToolbarItem, btnEditTimeToolbarItem)

	cnt := container.NewVSplit(timeToolbar, tt)
	cnt.Resize(fyne.NewSize(1000, 600))
	return cnt
}

func (av *AppView) AddTimeEntry() {
	// Add a time entry to current date
	loc, _ := time.LoadLocation("Europe/Berlin")
	today, err := av.repo.GetWorkday(time.Now().In(loc))
	if today == nil && err != nil {
		log.Debug("Create new workday")
		w := &db.Workday{
			Date: time.Now().In(loc),
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
	av.refreshAll()
}

func (av *AppView) deleteButtonFunc() {
	log.Info("Delete time entry", "item", av.selectedItem)
	// Check if is a day or a time entry selected
	if av.selectedItem == nil {
		dialog.ShowError(errors.New("no item selected"), av.window)
		return
	} else if av.selectedItem.Col == 0 {
		dialog.ShowConfirm("Delete Entry", "Do you really want to delete this workday?", func(b bool) {
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
					dialog.ShowError(errors.New("could not delete dataset"), av.window)
				}
			}
		}, av.window)
	} else {
		// Assure that there is an item selected
		wt, err := av.getTimeEntry(av.selectedItem)

		// Delete only if there is an existing item and no error
		// Else show an error dialog
		if wt != nil && err == nil {
			dialog.ShowConfirm("Delete Entry", "Do you really want to delete this time entry?", func(b bool) {
				if b {
					defer av.deleteTimeEntry(wt)

				} else {
					log.Debug("Delete time entry", "canceled", b)
				}
			}, av.window)
		} else {
			dialog.ShowError(errors.New("no item selected"), av.window)
		}
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
	wd, wdErr := av.repo.GetAllWorkday()
	if wdErr != nil {
		log.Error(wdErr)
	} else {
		av.workday = nil
		av.workday = wd
	}

	av.worktime = nil
	for i := 0; i < len(wd); i++ {
		wt, wtErr := av.repo.GetAllWorktime(wd[i])
		if wtErr != nil {
			log.Error(wtErr)
		} else {
			av.worktime = append(av.worktime, wt[0:]...)
		}
	}
}

func (av *AppView) getTimeEntry(item *widget.TableCellID) (*db.Worktime, error) {
	log.Info("Get time entry", "item", item)
	if item == nil {
		return nil, errors.New("no item selected")
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

	// Get the worktime by column (-1 because of the date column)
	if (av.selectedItem.Col - 1) < len(wtList) {
		return wtList[av.selectedItem.Col-1], nil
	} else {
		return nil, errors.New("no worktime found")
	}
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

/*
Combines the refresh of the data and the timetable
*/
func (av *AppView) refreshAll() {
	av.RefreshData()
	av.timetable.Refresh()
}

func (av *AppView) UnselectTableItem() {
	if av.selectedItem != nil {
		log.Debug("Unselect item", "item", av.selectedItem)
		av.timetable.Unselect(*av.selectedItem)
	} else {
		dialog.ShowError(errors.New("no item selected"), av.window)
	}
}

func (av *AppView) editButtonFunc() {
	loc, _ := time.LoadLocation("Europe/Berlin")

	log.Info("Edit time entry", "item", av.selectedItem)
	// Check if is a day or a time entry selected
	// Assure that there is an item selected
	if av.selectedItem == nil || av.selectedItem.Col == 0 {
		dialog.ShowError(errors.New("no time selected"), av.window)
		return
	}
	wt, err := av.getTimeEntry(av.selectedItem)

	if wt != nil && err == nil {
		timeEntry := widget.NewEntry()
		timeEntry.SetText(time.Now().In(loc).Format(time.TimeOnly))

		typeEntry := widget.NewSelectEntry([]string{"Begin", "End"})
		typeEntry.SetText(wt.Type)

		form := []*widget.FormItem{
			{Text: "Time (hh:mm:ss)", Widget: timeEntry, HintText: "Old entry: " + wt.Time.Format(time.TimeOnly)},
			{Text: "Type", Widget: typeEntry, HintText: "Old entry: " + wt.Type},
		}
		// Edit only if there is an existing item and no error
		log.Debug("Edit time entry", "worktime", wt)
		dia := dialog.NewForm("Edit Time Entry", "Edit", "Cancel", form, func(b bool) {
			if b {
				log.Debug("Edit time entry", "confirmed", b, "worktime", wt)
				// Parse the time
				timeEntry := form[0].Widget.(*widget.Entry).Text

				if timeEntry == "" {
					dialog.ShowError(errors.New("no time entry"), av.window)
					return
				}

				// Add the current date to the updated time entry
				tempTime := time.Now().In(loc).Format(time.DateOnly) + " " + timeEntry
				nt, err := time.ParseInLocation(time.DateTime, tempTime, loc)
				if err != nil {
					dialog.ShowError(err, av.window)
					return
				}
				log.Info("New time", "time", nt)

				prevCol := widget.TableCellID{Row: av.selectedItem.Row, Col: av.selectedItem.Col - 1} // -1 because of the date column
				if prevCol.Col >= 1 {                                                                 // If there is a previous entry
					prevTime, err := av.getTimeEntry(&prevCol)
					log.Debug("Previous time", "time", prevTime.Time)
					if err != nil {
						dialog.ShowError(err, av.window)
						return
					} else {
						// The edit shouldn't be before the previous entry
						log.Debug("Compare new-time with prev-time", "new-time", nt, "prev-time", prevTime.Time)
						if nt.Before(prevTime.Time) {
							dialog.ShowError(errors.New("time is before previous entry"), av.window)
							return
						}
					}
				}

				// The edit shouldn't be in future, it would be faking and does not make sense
				ct := time.Now().In(loc)
				log.Debug("Compare new-time with current-time", "new-time", nt, "current-time", ct)
				if ct.Before(nt) {
					dialog.ShowError(errors.New("time is in the future, do not fake 😉"), av.window)
					return
				}

				// Update the time entry
				wt.Time = nt
				wt.Type = form[1].Widget.(*widget.SelectEntry).Text

				_, errUpdate := av.repo.UpdateWorktime(wt)
				if errUpdate != nil {
					dialog.ShowError(err, av.window)
				} else {
					// Refresh *all data*
					av.refreshAll()
				}
			}
		}, av.window)
		dia.Resize(fyne.NewSize(400, 200))
		dia.Show()
	} else {
		dialog.ShowError(errors.New("no item selected"), av.window)
	}
}
