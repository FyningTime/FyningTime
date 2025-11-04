package view

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/FyningTime/FyningTime/app/model"
	"github.com/FyningTime/FyningTime/app/model/db"
	fwidget "github.com/FyningTime/FyningTime/app/widget"
)

type CalenderView struct {
	instruction *widget.Label
	dateChosen  *widget.Label

	container *fyne.Container

	calender *fwidget.Calendar
}

func (c *CalenderView) OnSelected(t time.Time) {
	loc, _ := time.LoadLocation("Europe/Berlin")

	c.instruction.SetText("Date selected:")
	c.dateChosen.SetText(t.In(loc).Format(model.DATEFORMAT))
}

func CreateCalendarView(w fyne.Window, vacations []*db.Vacation, selectedTime time.Time) *CalenderView {
	i := widget.NewLabel("Select a date")
	i.Alignment = fyne.TextAlignCenter
	l := widget.NewLabel("")
	l.Alignment = fyne.TextAlignCenter
	c := &CalenderView{instruction: i, dateChosen: l}

	xcalendar := fwidget.NewCalendar(w, vacations, selectedTime, c.OnSelected)
	content := container.NewBorder(l, nil, nil, nil, xcalendar)
	c.calender = xcalendar
	c.container = content
	return c
}

func (c *CalenderView) UpdateVacations(vacations []*db.Vacation) {
	if vacations != nil {
		c.calender.UpdateVacations(vacations)
	} else {
		c.calender.UpdateVacations([]*db.Vacation{})
	}
	c.calender.Refresh()
	c.container.Refresh()
}
