package main

import (
	"database/sql"
	"flag"

	"github.com/FyningTime/FyningTime/app/service"
	"github.com/FyningTime/FyningTime/app/view"

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
	db := GetDB("fyningtime.db")

	progName := "FyningTime"
	log.Info("Welcome to " + progName)

	// Test settings
	settings, err := service.ReadSettings()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Settings: " + settings.SavedPath)

	a := app.NewWithID("com.github.fyningtime.fyningtime")
	w := a.NewWindow(progName)

	av := new(view.AppView)
	av.CreateRepository(db)
	// INFO Date is the basic header, for each entry a new header is added
	av.AddHeaders([]string{"Date"})
	av.RefreshData()
	mv := av.CreateUI(w)

	setShortcuts(w, av.UnselectTableItem)

	// TODO implement me usefully!
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("FyneTime",
			fyne.NewMenuItem("Without function yet", func() {
				log.Info("Tapped show")
			}),
			fyne.NewMenuItem("About", func() {
				dialog.ShowInformation("About",
					"FyneTime is a simple time tracking application", w)
			}),
		)

		desk.SetSystemTrayMenu(m)
	}

	w.SetMainMenu(
		fyne.NewMainMenu(
			fyne.NewMenu("File",
				fyne.NewMenuItem("About", func() {
					log.Info("Tapped about")
					dialog.ShowInformation("About",
						"FyneTime is a simple time tracking application", w)
				}),
				fyne.NewMenuItem("Settings", func() {
					log.Info("Tapped settings")
				})),
		),
	)

	w.SetOnClosed(func() {
		log.Info("Closing database")
		db.Close()
		a.Quit()
	})
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

func setShortcuts(w fyne.Window, callbackShiftDelete func()) {
	// Shortcut for shift+del
	log.Info("Setting up shortcuts")
	deleteShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyDelete,
		Modifier: fyne.KeyModifierShift}
	w.Canvas().AddShortcut(deleteShortcut, func(shortcut fyne.Shortcut) {
		log.Info("Tapped shift+del")
		callbackShiftDelete()
	})
}
