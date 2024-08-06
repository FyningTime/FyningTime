package view

import (
	"database/sql"
	"main/app/model/db"
	"main/app/repo"
	"strconv"
	"time"

	"github.com/charmbracelet/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type AppView struct {
	// timetable *widget.Table represents a table widget used in the AppView struct.
	timetable *widget.Table

	// TODO add toolbar widget

	// Represents the headers for the timetable
	headers []string

	// Actual db abstraction
	worktime []*db.Worktime
	workday  []*db.Workday

	// Holds the database connection
	repo *repo.SQLiteRepository
}

func (av *AppView) CreateUI() *widget.Table {
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
					label.SetText(currentWt.Time.Local().Format(time.TimeOnly) + workType)
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
	av.timetable = tt
	return tt
}

func (av *AppView) AddTimeEntry() {
	// Add a time entry to current date
	today, err := av.repo.GetWorkday(time.Now())
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

			wtType := "Begin"
			if len(av.worktime)%2 != 0 {
				wtType = "End"
			}
			wt := &db.Worktime{
				Type:    wtType,
				Time:    time.Now(),
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
			wtType := "Begin"
			if len(allWt)%2 != 0 {
				wtType = "End"
			}

			wt := &db.Worktime{
				Type:    wtType,
				Time:    time.Now(),
				Workday: *today,
			}
			av.repo.AddWorktime(wt)
		}
	}

	// Refresh *all data*
	av.InitData()
	av.timetable.Refresh()
}

func (av *AppView) AddHeaders(headers []string) {
	av.headers = headers
}

func (av *AppView) AddRepository(db *sql.DB) {
	av.repo = repo.NewSQLiteRepository(db)

	// On start we also try to migrate the database
	log.Info("Migrating database")
	err := av.repo.Migrate()
	if err != nil {
		log.Fatal(err)
	}
}

// TODO maybe limit this for a specific date range like month
func (av *AppView) InitData() {
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
