package tui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/ritarock/passvault/domain"
	"github.com/rivo/tview"
)

type PasswordOptionsDialog struct {
	form       *tview.Form
	modal      *tview.Flex
	options    domain.PasswordOptions
	onGenerate func(domain.PasswordOptions)
	onCancel   func()
}

func NewPasswordOptionsDialog(onGenerate func(domain.PasswordOptions), onCancel func()) *PasswordOptionsDialog {
	pod := &PasswordOptionsDialog{
		form:       tview.NewForm(),
		options:    domain.DefaultPasswordOptions(),
		onGenerate: onGenerate,
		onCancel:   onCancel,
	}

	pod.setupForm()
	pod.setupModal()

	return pod
}

func (pod *PasswordOptionsDialog) setupForm() {
	pod.form.SetTitle(" Password Options ").
		SetBorder(true).
		SetBorderColor(ColorPrimary)

	pod.form.AddInputField("Length (8-64)", strconv.Itoa(pod.options.Length), 10, func(text string, lastChar rune) bool {
		if lastChar < '0' || lastChar > '9' {
			return false
		}
		return len(text) <= 2
	}, func(text string) {
		if length, err := strconv.Atoi(text); err == nil {
			pod.options.Length = length
		}
	})

	pod.form.AddCheckbox("Include lowercase (a-z)", pod.options.IncludeLowercase, func(checked bool) {
		pod.options.IncludeLowercase = checked
	})

	pod.form.AddCheckbox("Include uppercase (A-Z)", pod.options.IncludeUppercase, func(checked bool) {
		pod.options.IncludeUppercase = checked
	})

	pod.form.AddCheckbox("Include digits (0-9)", pod.options.IncludeDigits, func(checked bool) {
		pod.options.IncludeDigits = checked
	})

	pod.form.AddCheckbox("Include symbols (!@#$...)", pod.options.IncludeSymbols, func(checked bool) {
		pod.options.IncludeSymbols = checked
	})

	pod.form.AddButton("Generate", func() {
		pod.onGenerate(pod.options)
	})

	pod.form.AddButton("Cancel", func() {
		pod.onCancel()
	})

	pod.form.SetButtonsAlign(tview.AlignCenter)

	pod.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pod.onCancel()
			return nil
		}
		return event
	})
}

func (pod *PasswordOptionsDialog) setupModal() {
	pod.modal = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(pod.form, 16, 0, true).
			AddItem(nil, 0, 1, false), 45, 0, true).
		AddItem(nil, 0, 1, false)
}

func (pod *PasswordOptionsDialog) GetPrimitive() tview.Primitive {
	return pod.modal
}
