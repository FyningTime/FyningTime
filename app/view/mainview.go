package view

import (
	"database/sql"
	"main/app/model"
	"main/app/model/db"
	"main/app/repo"
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

	// Represents the data entry for the timetable
	timeEntries []*model.TimeEntry

	// Actual db abstraction
	worktime []*db.Worktime
	workday  *db.Workday

	// Holds the database connection
	repo *repo.SQLiteRepository
}

func (av *AppView) CreateUI() *widget.Table {
	tt := widget.NewTable(
		func() (int, int) {
			return len(av.timeEntries) + 1, len(av.headers)
		},
		func() fyne.CanvasObject {
			// Defines how wide the cell is
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)

			// Header-Zeile
			if i.Row == 0 {
				label.SetText(av.headers[i.Col])
			} else {
				// Datenzeilen
				entry := av.timeEntries[i.Row-1]
				switch i.Col {
				case 0:
					label.SetText(entry.DATE.Format(time.DateOnly))
				default:
					entriesSize := len(entry.ENTRIES)
					if entriesSize >= i.Col {
						if entriesSize%2 != 0 && i.Col%2 != 0 {
							label.Importance = widget.HighImportance
						}
						label.SetText(entry.ENTRIES[i.Col-1].Format(time.TimeOnly))
					}
				}
			}
		},
	)
	//tt.Resize(fyne.NewSize(1000, 400))
	tt.SetColumnWidth(1, 100)
	av.timetable = tt
	return tt
}

func (av *AppView) AddTimeEntry() {
	// Add a time entry to current date
	var today *model.TimeEntry

	// Check if for the current date exist an entry
	entries := av.timeEntries
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if entry.DATE.Day() == time.Now().Day() {
			today = entry.AddTime(time.Now())
			entrySize := len(entry.ENTRIES)
			log.Debug("Entry:", "size", entrySize)
			entryHeader := "Begin"
			if entrySize%2 == 0 {
				entryHeader = "End"
			}
			av.headers = append(av.headers, entryHeader)
			log.Info("Time updated:", "entry", entry)
			break
		}
	}

	if today == nil {
		te := model.New()
		te.ENTRIES = append(te.ENTRIES, time.Now())
		av.timeEntries = append(av.timeEntries, te)
		entrySize := len(te.ENTRIES)
		log.Debug("Entry", "size", entrySize)
		entryHeader := "Begin"
		if entrySize%2 == 0 {
			entryHeader = "End"
		}
		av.headers = append(av.headers, entryHeader)
		log.Info("Time added:", "time-entry", te)
	}

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
