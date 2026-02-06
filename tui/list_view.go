package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/ritarock/passvault/domain"
	"github.com/rivo/tview"
)

type ListView struct {
	app             *App
	container       *tview.Flex
	table           *tview.Table
	help            *tview.TextView
	searchField     *tview.InputField
	entries         []*domain.Entry
	filteredEntries []*domain.Entry
}

func NewListView(app *App) *ListView {
	lv := &ListView{
		app:         app,
		table:       tview.NewTable(),
		help:        tview.NewTextView(),
		searchField: tview.NewInputField(),
	}

	lv.setupTable()
	lv.setupSearchField()
	lv.setupHelp()
	lv.setupContainer()

	return lv
}

func (lv *ListView) setupTable() {
	lv.table.SetBorder(true).
		SetTitle(" PassVault ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(ColorPrimary)

	lv.table.SetSelectable(true, false)
	lv.table.SetSelectedStyle(tcell.StyleDefault.
		Background(ColorPrimary).
		Foreground(tcell.ColorWhite))

	lv.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			lv.app.ShowForm("")
			return nil
		case 'd':
			lv.deleteSelected()
			return nil
		case 'q':
			lv.app.Stop()
			return nil
		case '/':
			lv.app.app.SetFocus(lv.searchField)
			return nil
		}

		switch event.Key() {
		case tcell.KeyEnter:
			lv.viewSelected()
			return nil
		}

		return event
	})
}

func (lv *ListView) setupSearchField() {
	lv.searchField.SetLabel(" Search: ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetLabelColor(ColorPrimary).
		SetFieldTextColor(tcell.ColorWhite)

	lv.searchField.SetChangedFunc(func(text string) {
		lv.filterEntries(text)
		lv.renderTable()
	})

	lv.searchField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			lv.app.app.SetFocus(lv.table)
			return nil
		}
		return event
	})
}

func (lv *ListView) setupHelp() {
	lv.help.SetText("[/] Search  [a] Add  [Enter] View  [d] Delete  [q] Quit").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(ColorSecondary)
}

func (lv *ListView) setupContainer() {
	lv.container = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(lv.searchField, 1, 0, false).
		AddItem(lv.table, 0, 1, true).
		AddItem(lv.help, 1, 0, false)
}

func (lv *ListView) GetPrimitive() tview.Primitive {
	return lv.container
}

func (lv *ListView) Refresh() {
	entries, err := lv.app.listEntriesUc.Execute()
	if err != nil {
		lv.app.ShowError(fmt.Sprintf("Failed to load entries: %v", err))
		return
	}

	lv.entries = entries
	lv.searchField.SetText("")
	lv.filterEntries("")

	lv.renderTable()
}

func (lv *ListView) filterEntries(query string) {
	if query == "" {
		lv.filteredEntries = lv.entries
		return
	}

	query = strings.ToLower(query)
	lv.filteredEntries = nil
	for _, entry := range lv.entries {
		if strings.Contains(strings.ToLower(entry.Title), query) ||
			strings.Contains(strings.ToLower(entry.Username), query) ||
			strings.Contains(strings.ToLower(entry.URL), query) {
			lv.filteredEntries = append(lv.filteredEntries, entry)
		}
	}
}

func (lv *ListView) renderTable() {
	lv.table.Clear()

	lv.table.SetCell(0, 0, tview.NewTableCell("Title").
		SetTextColor(ColorPrimary).
		SetSelectable(false))
	lv.table.SetCell(0, 1, tview.NewTableCell("Username").
		SetTextColor(ColorPrimary).
		SetSelectable(false))
	lv.table.SetCell(0, 2, tview.NewTableCell("URL").
		SetTextColor(ColorPrimary).
		SetSelectable(false))
	lv.table.SetCell(0, 3, tview.NewTableCell("Notes").
		SetTextColor(ColorPrimary).
		SetSelectable(false))
	lv.table.SetCell(0, 4, tview.NewTableCell("CreatedAt").
		SetTextColor(ColorPrimary).
		SetSelectable(false))

	for i, entry := range lv.filteredEntries {
		row := i + 1
		lv.table.SetCell(row, 0, tview.NewTableCell(entry.Title).
			SetTextColor(tcell.ColorWhite))
		lv.table.SetCell(row, 1, tview.NewTableCell(entry.Username).
			SetTextColor(ColorSecondary))
		lv.table.SetCell(row, 2, tview.NewTableCell(entry.URL).
			SetTextColor(ColorSecondary))
		lv.table.SetCell(row, 3, tview.NewTableCell(entry.Notes).
			SetTextColor(ColorSecondary))
		lv.table.SetCell(row, 4, tview.NewTableCell(entry.CreatedAt.Format("2006-01-02 15:04")).
			SetTextColor(ColorSecondary))
	}

	if len(lv.filteredEntries) == 0 {
		message := "No entries yet. Press 'a' to add one."
		if len(lv.entries) > 0 {
			message = "No matching entries found."
		}
		lv.table.SetCell(1, 0, tview.NewTableCell(message).
			SetTextColor(ColorSecondary).
			SetAlign(tview.AlignCenter))
	}

	lv.table.Select(1, 0)
}

func (lv *ListView) viewSelected() {
	if len(lv.filteredEntries) == 0 {
		return
	}

	row, _ := lv.table.GetSelection()
	if row < 1 {
		return
	}

	index := row - 1
	if index >= len(lv.filteredEntries) {
		return
	}

	lv.app.ShowDetail(lv.filteredEntries[index].ID)
}

func (lv *ListView) deleteSelected() {
	if len(lv.filteredEntries) == 0 {
		return
	}

	row, _ := lv.table.GetSelection()
	if row < 1 {
		return
	}

	index := row - 1
	if index >= len(lv.filteredEntries) {
		return
	}

	entry := lv.filteredEntries[index]
	lv.app.ShowConfirm(
		fmt.Sprintf("Delete '%s'?", entry.Title),
		func() {
			if err := lv.app.deleteEntryUc.Execute(entry.ID); err != nil {
				lv.app.ShowError(fmt.Sprintf("Failed to delete entry: %v", err))
				return
			}
			lv.Refresh()
		},
	)
}
