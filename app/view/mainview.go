package view

import (
	"fmt"
	"main/app/model"
	"strconv"
	"time"

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
}

func (av *AppView) CreateUI() *widget.Table {
	tt := widget.NewTable(
		func() (int, int) {
			return len(av.timeEntries) + 1, len(av.headers)
		},
		func() fyne.CanvasObject {
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
					entriesSize := len(entry.ENTRY)
					if entriesSize >= i.Col {
						label.SetText(entry.ENTRY[i.Col-1].Format(time.TimeOnly))
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
			entrySize := len(entry.ENTRY)
			fmt.Println("Entry size: ", entrySize)
			av.headers = append(av.headers, "Entry #"+strconv.Itoa(entrySize))
			fmt.Printf("Time updated %v\n", entry)
			break
		}
	}

	if today == nil {
		te := model.New()
		te.ENTRY = append(te.ENTRY, time.Now())
		av.timeEntries = append(av.timeEntries, te)
		entrySize := len(te.ENTRY)
		fmt.Println("Entry size: ", entrySize)
		av.headers = append(av.headers, "Entry #"+strconv.Itoa(entrySize))
		fmt.Printf("Time added %v\n", te)
	}

	av.timetable.Refresh()
}

func (av *AppView) AddHeaders(headers []string) {
	av.headers = headers
}
