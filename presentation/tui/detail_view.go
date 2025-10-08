package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/ritarock/passvault/domain/entity"
	"github.com/rivo/tview"
)

type DetailView struct {
	app       *App
	container *tview.Flex
	textView  *tview.TextView
	help      *tview.TextView
	entry     *entity.Entry
}

func NewDetailView(app *App) *DetailView {
	dv := &DetailView{
		app:      app,
		textView: tview.NewTextView(),
		help:     tview.NewTextView(),
	}

	dv.setupTextView()
	dv.setupHelp()
	dv.setupContainer()

	return dv
}

func (dv *DetailView) setupTextView() {
	dv.textView.SetDynamicColors(true).
		SetBorder(true).
		SetTitle(" Entry Details ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(ColorPrimary)

	dv.textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'e':
			if dv.entry != nil {
				dv.app.ShowForm(dv.entry.ID)
			}
			return nil
		case 'c':
			dv.copyPassword()
			return nil
		}

		switch event.Key() {
		case tcell.KeyEscape:
			dv.app.ShowList()
			return nil
		}

		return event
	})
}

func (dv *DetailView) setupHelp() {
	dv.help.SetText("[e] Edit  [c] Copy Password  [ESC] Back").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(ColorSecondary)
}

func (dv *DetailView) setupContainer() {
	dv.container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(dv.textView, 0, 1, true).
		AddItem(dv.help, 1, 0, false)
}

func (dv *DetailView) GetPrimitive() tview.Primitive {
	return dv.container
}

func (dv *DetailView) SetEntry(id string) {
	entry, err := dv.app.getEntryUc.Execute(id)
	if err != nil {
		dv.app.ShowError(fmt.Sprintf("Failed to load entry: %v", err))
		dv.app.ShowList()
		return
	}

	dv.entry = entry
	dv.render()
}

func (dv *DetailView) render() {
	if dv.entry == nil {
		return
	}

	var content strings.Builder

	content.WriteString(fmt.Sprintf("[::b]Name:[-:-:-]\n%s\n\n", dv.entry.Name))
	content.WriteString(fmt.Sprintf("[::b]Password:[-:-:-]\n%s\n\n", maskPassword(dv.entry.Password)))

	if dv.entry.URL != "" {
		content.WriteString(fmt.Sprintf("[::b]URL:[-:-:-]\n%s\n\n", dv.entry.URL))
	}

	if dv.entry.Notes != "" {
		content.WriteString(fmt.Sprintf("[::b]Notes:[-:-:-]\n%s\n\n", dv.entry.Notes))
	}

	content.WriteString(fmt.Sprintf("[::b]Created:[-:-:-] %s\n", dv.entry.CreatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("[::b]Updated:[-:-:-] %s\n", dv.entry.UpdatedAt.Format("2006-01-02 15:04:05")))

	dv.textView.SetText(content.String())
}

func (dv *DetailView) copyPassword() {
	if dv.entry == nil {
		return
	}

	if err := clipboard.WriteAll(dv.entry.Password); err != nil {
		dv.app.ShowError(fmt.Sprintf("Failed to copy password: %v", err))
		return
	}

	modal := tview.NewModal().
		SetText("Password copied to clipboard!").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			dv.app.pages.RemovePage("copied")
		})
	modal.SetBackgroundColor(tcell.ColorDefault)
	modal.SetBorderColor(ColorSuccess)
	dv.app.pages.AddPage("copied", modal, true, true)
}

func maskPassword(password string) string {
	return strings.Repeat("*", len(password))
}
