package main

import (
	"main/app/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.github.fynetime.fynetime")
	// FIXME: SetIcon is not available in the current version of Fyne
	//a.SetIcon(fyne.NewStaticResource("logo.png", nil))
	w := a.NewWindow("FyneTime")

	//toolbar := view.toolbar.makeUI()
	// TODO get time entries from the database and pass them to the view
	av := new(view.AppView)
	av.AddHeaders([]string{"Date"})
	mv := av.CreateUI()
	//av.AddTimeEntry()

	btn := widget.NewButton("Add Time Entry", func() {
		av.AddTimeEntry()
	})

	w.Resize(fyne.NewSize(1000, 500))
	w.SetContent(container.NewVSplit(btn, mv))
	w.ShowAndRun()
}
