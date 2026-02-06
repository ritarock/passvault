package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/ritarock/passvault/domain"
	"github.com/rivo/tview"
)

type FormView struct {
	app       *App
	container *tview.Flex
	form      *tview.Form
	help      *tview.TextView
	entryID   string
	isEdit    bool
}

func NewFormView(app *App) *FormView {
	fv := &FormView{
		app:  app,
		form: tview.NewForm(),
		help: tview.NewTextView(),
	}

	fv.setupHelp()
	fv.setupContainer()

	return fv
}

func (fv *FormView) setupHelp() {
	fv.help.SetText("[Tab] Next Field  [Shift+Tab] Previous Field  [Enter] Generate/Save  [ESC] Cancel").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(ColorSecondary)
}

func (fv *FormView) setupContainer() {
	fv.container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(fv.form, 0, 1, true).
		AddItem(fv.help, 1, 0, false)
}

func (fv *FormView) GetPrimitive() tview.Primitive {
	return fv.container
}

func (fv *FormView) SetEntry(id string) {
	fv.entryID = id
	fv.isEdit = id != ""

	fv.form.Clear(true)

	if fv.isEdit {
		fv.form.SetTitle(" Edit Entry ").SetBorder(true).SetBorderColor(ColorPrimary)
		fv.loadEntry()
	} else {
		fv.form.SetTitle(" Add Entry ").SetBorder(true).SetBorderColor(ColorPrimary)
		fv.setupNewForm()
	}
}

func (fv *FormView) loadEntry() {
	entry, err := fv.app.getEntryUc.Execute(fv.entryID)
	if err != nil {
		fv.app.ShowError(fmt.Sprintf("Failed to load entry: %v", err))
		fv.app.ShowList()
		return
	}

	fv.setupFormFields(entry.Title, entry.Username, entry.Password, entry.URL, entry.Notes)
}

func (fv *FormView) setupNewForm() {
	fv.setupFormFields("", "", "", "", "")
}

func (fv *FormView) setupFormFields(title, username, password, url, notes string) {
	fv.form.AddInputField("Title", title, 40, nil, nil)
	fv.form.AddInputField("Username", username, 40, nil, nil)
	fv.form.AddPasswordField("Password", password, 40, '*', nil)
	fv.form.AddInputField("URL", url, 40, nil, nil)
	fv.form.AddTextArea("Notes", notes, 40, 3, 0, nil)

	fv.form.AddButton("Generate Password", fv.generatePassword)
	fv.form.AddButton("Save", fv.save)
	fv.form.AddButton("Cancel", func() {
		fv.app.ShowList()
	})

	fv.form.SetButtonsAlign(tview.AlignCenter)
	fv.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			fv.app.ShowList()
			return nil
		}
		return event
	})
}

func (fv *FormView) save() {
	title := fv.form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
	username := fv.form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	password := fv.form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	url := fv.form.GetFormItemByLabel("URL").(*tview.InputField).GetText()
	notes := fv.form.GetFormItemByLabel("Notes").(*tview.TextArea).GetText()

	if title == "" {
		fv.app.ShowError("Title is required")
		return
	}
	if password == "" {
		fv.app.ShowError("Password is required")
		return
	}

	var err error
	if fv.isEdit {
		err = fv.app.updateEntryUc.Execute(fv.entryID, title, username, password, url, notes)
	} else {
		err = fv.app.createEntryUc.Execute(title, username, password, url, notes)
	}

	if err != nil {
		fv.app.ShowError(fmt.Sprintf("Failed to save entry: %v", err))
		return
	}

	fv.app.ShowList()
}

func (fv *FormView) generatePassword() {
	fv.app.ShowPasswordOptionsDialog(func(opts domain.PasswordOptions) {
		password, err := fv.app.passwordGen.GenerateWithOptions(opts)
		if err != nil {
			fv.app.ShowError(fmt.Sprintf("Failed to generate password: %v", err))
			return
		}

		passwordField := fv.form.GetFormItemByLabel("Password").(*tview.InputField)
		passwordField.SetText(password)
	})
}
