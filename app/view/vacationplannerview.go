package view

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/FyningTime/FyningTime/app/model"
	"github.com/FyningTime/FyningTime/app/model/db"
	"github.com/FyningTime/FyningTime/app/repo"
	"github.com/charmbracelet/log"
	datepicker "github.com/sdassow/fyne-datepicker"
)

type VacationPlannerView struct {
	// Business logic
	vacations []*db.Vacation
	repo      *repo.SQLiteRepository

	// UI
	av               *AppView
	container        *fyne.Container
	selectedVacation widget.ListItemID
}

func NewVacationPlannerView(
	av *AppView,
	repo *repo.SQLiteRepository,
	vacations []*db.Vacation,
) *VacationPlannerView {
	return &VacationPlannerView{
		vacations: vacations,
		repo:      repo,
		av:        av,
	}
}

func CreateVacationPlannerView(av *AppView, repo *repo.SQLiteRepository, vacations []*db.Vacation) *VacationPlannerView {
	vpv := NewVacationPlannerView(av, repo, vacations)
	log.Debug("CreateVacationPlannerView", "VacationPlannerView", vpv)

	vl := widget.NewList(
		func() int {
			return len(vpv.vacations)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Vacations")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(
				vpv.vacations[i].StartDate.Format(model.DATEFORMAT) + " - " + vpv.vacations[i].EndDate.Format(model.DATEFORMAT),
			)
		},
	)

	vl.OnSelected = func(id widget.ListItemID) {
		vpv.selectedVacation = id
	}

	header := widget.NewLabel("Vacations Planner")
	header.Alignment = fyne.TextAlignCenter
	header.TextStyle = fyne.TextStyle{Bold: true}

	btnAddVacationToolbarItem := widget.NewToolbarAction(theme.ContentAddIcon(), vpv.addVacationForm)
	btnDeleteTimeToolbarItem := widget.NewToolbarAction(theme.ContentRemoveIcon(), vpv.deleteVacationForm)
	btnEditTimeToolbarItem := widget.NewToolbarAction(theme.DocumentIcon(), vpv.editVacationForm)

	toolbar := widget.NewToolbar(
		btnAddVacationToolbarItem,
		btnDeleteTimeToolbarItem,
		btnEditTimeToolbarItem,
	)

	c := container.NewBorder(toolbar, nil, nil, nil, vl)
	vpv.container = c
	return vpv
}

func (vpv *VacationPlannerView) addVacationForm() {
	startDate := widget.NewEntry()
	startDate.SetPlaceHolder("01.01.1970")
	startDate.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when := time.Now()

		picker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				startDate.SetText(when.Format(model.DATEFORMAT))
			}
		})

		dialog.ShowCustomConfirm(
			"Select a date",
			"OK",
			"Cancel",
			picker,
			picker.OnActioned,
			vpv.av.window,
		)
	})

	endDate := widget.NewEntry()
	endDate.SetPlaceHolder("01.01.1970")
	endDate.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when := time.Now()

		picker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				endDate.SetText(when.Format(model.DATEFORMAT))
			}
		})

		dialog.ShowCustomConfirm(
			"Select a date",
			"OK",
			"Cancel",
			picker,
			picker.OnActioned,
			vpv.av.window,
		)
	})

	dialog.ShowForm("Add Vacation", "Add", "Cancel",
		[]*widget.FormItem{
			{Text: "Start Date", Widget: startDate},
			{Text: "End Date", Widget: endDate},
		}, func(submitted bool) {
			if submitted {
				s, err := time.Parse(model.DATEFORMAT, startDate.Text)
				if err != nil {
					log.Error(err)
					dialog.ShowError(err, vpv.av.window)
				}

				e, err := time.Parse(model.DATEFORMAT, endDate.Text)
				if err != nil {
					log.Error(err)
					dialog.ShowError(err, vpv.av.window)
				}

				vpv.addVacation(&db.Vacation{
					StartDate: s,
					EndDate:   e,
				})
			}
		}, vpv.av.window)
}

func (vpv *VacationPlannerView) deleteVacationForm() {
	// TODO: Implement delete vacation form
}

func (vpv *VacationPlannerView) editVacationForm() {
	// TODO: Implement edit vacation form
}

func (vpv *VacationPlannerView) UpdateVacations(vacations []*db.Vacation) {
	vpv.vacations = vacations
	vpv.container.Refresh()
}

func (vpv *VacationPlannerView) addVacation(v *db.Vacation) {
	_, err := vpv.repo.AddVacation(v)
	if err != nil {
		log.Error(err)
		dialog.ShowError(err, vpv.av.window)
	} else {
		vpv.av.RefreshData()
	}
}
