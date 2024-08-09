package main

import (
	"database/sql"
	"flag"
	"main/app/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var devFlag bool
	initLogging(devFlag)

	progName := "FyningTime"
	log.Info("Welcome to " + progName)

	a := app.NewWithID("com.github.fyningtime.fyningtime")
	w := a.NewWindow(progName)

	//toolbar := view.toolbar.makeUI()
	av := new(view.AppView)
	av.AddRepository(GetDB("fyningtime.db"))
	// INFO Date is the basic header, for each entry a new header is added
	av.AddHeaders([]string{"Date"})
	av.InitData()
	mv := av.CreateUI()
	mv.OnSelected = func(id widget.TableCellID) {
		log.Debug("Selected", "id", id)
	}
	//av.AddTimeEntry()

	// TODO move btn to main view
	btn := widget.NewButton("Add Time Entry", func() {
		av.AddTimeEntry()
	})

	w.SetContent(container.NewVSplit(btn, mv))
	w.Resize(fyne.NewSize(800, 600))
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

func initLogging(devFlag bool) {
	flag.BoolVar(&devFlag, "d", false, "Development flag")
	flag.Parse()
	if devFlag {
		log.Info("Development mode enabled")
		log.SetLevel(log.DebugLevel)
	}
}
