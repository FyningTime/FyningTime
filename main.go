package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/FyningTime/FyningTime/app/service"
	apptheme "github.com/FyningTime/FyningTime/app/theme"
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

	progName := "FyningTime"
	log.Info("Welcome to " + progName)
	a := app.NewWithID("com.github.fyningtime.fyningtime")

	// Read settings
	settings := service.ReadProperties(a)
	log.Debugf("Settings: %+v", settings)

	// Open database
	db := GetDB(settings.SavedDbPath)

	switch settings.ThemeVariant {
	case 1:
		a.Settings().SetTheme(apptheme.NewPastelleDark())

	case 2:
		a.Settings().SetTheme(apptheme.NewPastelleLight())

	default:
		a.Settings().SetTheme(apptheme.NewPastelleTheme())
	}

	w := a.NewWindow(progName)

	av := new(view.AppView)
	av.CreateRepository(db)
	av.SetBaseHeaders([]string{"Date", "Time", "Break", "Overtime"})

	mv := av.CreateUI(w, a)
	av.RefreshData()

	// Set shortcuts
	setShortcuts(w, CreateAppShortcuts(av))

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("FyneTime",
			fyne.NewMenuItem("Add current time", func() {
				av.AddTimeEntry()
			}),
			// Show overtime in dialog
			fyne.NewMenuItem("Show overtime", func() {
				overtime := av.GetOvertime()
				dialog.ShowInformation("Overtime", fmt.Sprintf("Current overtime: %s", overtime), w)
			}),
			fyne.NewMenuItem("About", func() {
				dialog.NewInformation("About",
					"FyningTime is a simple time tracking application", w)
			}),
		)

		desk.SetSystemTrayMenu(m)
	}

	w.SetMainMenu(
		fyne.NewMainMenu(
			fyne.NewMenu("File",
				fyne.NewMenuItem("About", func() {
					dialog.ShowInformation("About",
						"FyningTime is a simple time tracking application", w)
				}),
				fyne.NewMenuItem("Settings", func() {
					view.GetSettingsView(w, a).Show()
				})),
		),
	)

	w.SetOnClosed(func() {
		log.Info("Closing database")
		db.Close()
		a.Quit()
	})
	w.SetContent(mv)
	w.Resize(fyne.NewSize(700, 600))
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

func setShortcuts(w fyne.Window, shortcuts []AppShortcuts) {
	log.Info("Setting up shortcuts")

	canvas := w.Canvas()

	for _, sc := range shortcuts {
		customShortcut := &desktop.CustomShortcut{
			KeyName:  sc.keyname,
			Modifier: sc.modifier,
		}

		canvas.AddShortcut(customShortcut, func(shortcut fyne.Shortcut) {
			log.Infof("Tapped shortcut: %s", sc.name)
			sc.callback()
		})
	}
}

type AppShortcuts struct {
	// The shortcut name
	name     string
	callback func()

	keyname  fyne.KeyName
	modifier fyne.KeyModifier
}

func CreateAppShortcuts(av *view.AppView) []AppShortcuts {
	unselectShortcut := AppShortcuts{
		name:     "Unselect Table Item",
		callback: func() { av.UnselectTableItem() },
		keyname:  fyne.KeyU,
		modifier: fyne.KeyModifierControl,
	}

	deleteSelectedTimeEntryShortcut := AppShortcuts{
		name:     "Delete Selected Time Entry",
		callback: func() { av.DeleteSelectedTimeEntry() },
		keyname:  fyne.KeyDelete,
		modifier: fyne.KeyModifierControl,
	}

	addNewTimeEntryShortcut := AppShortcuts{
		name:     "Add New Time Entry",
		callback: func() { av.AddTimeEntry() },
		keyname:  fyne.KeyN,
		modifier: fyne.KeyModifierControl,
	}

	editSelectedTimeEntryShortcut := AppShortcuts{
		name:     "Edit Selected Time Entry",
		callback: func() { av.EditSelectedTimeEntry() },
		keyname:  fyne.KeyE,
		modifier: fyne.KeyModifierControl,
	}

	refreshDataShortcut := AppShortcuts{
		name:     "Refresh Data",
		callback: func() { av.RefreshData() },
		keyname:  fyne.KeyR,
		modifier: fyne.KeyModifierControl,
	}

	return []AppShortcuts{
		unselectShortcut,
		deleteSelectedTimeEntryShortcut,
		addNewTimeEntryShortcut,
		editSelectedTimeEntryShortcut,
		refreshDataShortcut,
	}
}
