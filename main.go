package main

import (
	"database/sql"
	"main/app/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	progName := "FyningTime"
	log.Info("Welcome to " + progName)

	a := app.NewWithID("com.github.fyningtime.fyningtime")
	w := a.NewWindow(progName)

	//toolbar := view.toolbar.makeUI()
	av := new(view.AppView)
	// TODO get time entries from the database and pass them to the view
	av.AddRepository(GetDB("fyningtime.db"))
	// INFO Date is the basic header, for each entry a new header is added
	av.AddHeaders([]string{"Date"})
	mv := av.CreateUI()
	//av.AddTimeEntry()

	// TODO move btn to main view
	btn := widget.NewButton("Add Time Entry", func() {
		av.AddTimeEntry()
	})

	w.Resize(fyne.NewSize(1000, 500))
	w.SetContent(container.NewVSplit(btn, mv))
	w.ShowAndRun()
}

func GetDB(filePath string) *sql.DB {
	log.Info("Opening database to: " + filePath)
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
