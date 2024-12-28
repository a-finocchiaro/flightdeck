package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Fr24Modal struct {
	Modal      tview.Primitive
	Form       *tview.Form
	ActionFunc func(int, string)
}

func NewFr24Modal(text string) *Fr24Modal {
	m := Fr24Modal{}
	modal := func(p tview.Primitive, width int, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	m.Form = tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor).
		AddInputField("Airport IATA:", "", 5, nil, nil)

	m.Form.SetBorder(true).SetTitle("Select Airport")

	m.AddButtons([]string{"Cancel", "Search"})
	m.Modal = modal(m.Form, 40, 10)

	return &m
}

func (m *Fr24Modal) SetActionFunc(fn func(int, string)) *Fr24Modal {
	m.ActionFunc = fn
	return m
}

// Adds buttons to the ModalForm
func (m *Fr24Modal) AddButtons(labels []string) *Fr24Modal {
	for index, label := range labels {
		func(i int, l string) {
			m.Form.AddButton(label, func() {
				if m.ActionFunc != nil {
					m.ActionFunc(i, l)
				}
			})
			button := m.Form.GetButton(m.Form.GetButtonCount() - 1)
			button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Key() {
				case tcell.KeyDown, tcell.KeyRight:
					return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
				case tcell.KeyUp, tcell.KeyLeft:
					return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
				}
				return event
			})
		}(index, label)
	}
	return m
}
