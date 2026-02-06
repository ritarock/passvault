package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/ritarock/passvault/domain"
	"github.com/ritarock/passvault/service"
	"github.com/rivo/tview"
)

type App struct {
	app           *tview.Application
	pages         *tview.Pages
	listView      *ListView
	detailView    *DetailView
	formView      *FormView
	listEntriesUc *service.ListEntriesUsecase
	getEntryUc    *service.GetEntryUsecase
	createEntryUc *service.CreateEntryUsecase
	updateEntryUc *service.UpdateEntryUsecase
	deleteEntryUc *service.DeleteEntryUsecase
	passwordGen   *domain.PasswordGenerator
}

func NewApp(
	listEntriesUc *service.ListEntriesUsecase,
	getEntryUc *service.GetEntryUsecase,
	createEntryUc *service.CreateEntryUsecase,
	updateEntryUc *service.UpdateEntryUsecase,
	deleteEntryUc *service.DeleteEntryUsecase,
) *App {
	app := &App{
		app:           tview.NewApplication(),
		pages:         tview.NewPages(),
		listEntriesUc: listEntriesUc,
		getEntryUc:    getEntryUc,
		createEntryUc: createEntryUc,
		updateEntryUc: updateEntryUc,
		deleteEntryUc: deleteEntryUc,
		passwordGen:   domain.NewPasswordGenerator(),
	}

	app.listView = NewListView(app)
	app.detailView = NewDetailView(app)
	app.formView = NewFormView(app)

	app.pages.AddPage("list", app.listView.GetPrimitive(), true, true)
	app.pages.AddPage("detail", app.detailView.GetPrimitive(), true, false)
	app.pages.AddPage("form", app.formView.GetPrimitive(), true, false)

	app.app.SetRoot(app.pages, true)

	return app
}

func (a *App) Run() error {
	return a.app.Run()
}

func (a *App) Stop() {
	a.app.Stop()
}

func (a *App) ShowList() {
	a.listView.Refresh()
	a.pages.SwitchToPage("list")
}

func (a *App) ShowDetail(id string) {
	a.detailView.SetEntry(id)
	a.pages.SwitchToPage("detail")
}

func (a *App) ShowForm(id string) {
	a.formView.SetEntry(id)
	a.pages.SwitchToPage("form")
}

func (a *App) ShowError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("error")
		})
	modal.SetBackgroundColor(tcell.ColorDefault)
	modal.SetBorderColor(ColorDanger)
	a.pages.AddPage("error", modal, true, true)
}

func (a *App) ShowConfirm(message string, onConfirm func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("confirm")
			if buttonIndex == 0 {
				onConfirm()
			}
		})
	modal.SetBackgroundColor(tcell.ColorDefault)
	modal.SetBorderColor(ColorPrimary)
	a.pages.AddPage("confirm", modal, true, true)
}

func (a *App) ShowPasswordOptionsDialog(onGenerate func(domain.PasswordOptions)) {
	dialog := NewPasswordOptionsDialog(
		func(opts domain.PasswordOptions) {
			a.pages.RemovePage("password-options")
			onGenerate(opts)
		},
		func() {
			a.pages.RemovePage("password-options")
		},
	)
	a.pages.AddPage("password-options", dialog.GetPrimitive(), true, true)
}
