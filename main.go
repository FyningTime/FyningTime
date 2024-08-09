package main

import (
	"database/sql"
	"flag"
	"main/app/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
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
	av.RefreshData()
	mv := av.CreateUI(w)
	//av.AddTimeEntry()

	// TODO implement me usefully!
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("FyneTime",
			fyne.NewMenuItem("Without function yet", func() {
				log.Info("Tapped show")
			}))
		desk.SetSystemTrayMenu(m)
	}

	w.SetContent(mv)
	w.Resize(fyne.NewSize(1000, 600))
	w.ShowAndRun()
}

func GetDB(filePath string) *sql.DB {
	log.Info("Opening database to: " + filePath)
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
		dialog.ShowError(err, nil)
		return nil
	}
	// Fremdschlüsselunterstützung aktivieren
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
		dialog.ShowError(err, nil)
		return nil
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
