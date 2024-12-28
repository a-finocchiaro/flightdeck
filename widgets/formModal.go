package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormModal struct {
	Modal      tview.Primitive
	Form       *tview.Form
	ActionFunc func(int, string)
}

type InputFields struct {
	Label       string
	Placeholder string
	Length      int
	Accept      func(string, rune) bool
	OnChange    func(string)
}

func NewFormModal(buttons []string, inputFields []InputFields) *FormModal {
	m := FormModal{}
	modal := func(p tview.Primitive, width int, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	m.Form = tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor)

	for _, field := range inputFields {
		m.Form.AddInputField(
			field.Label,
			field.Placeholder,
			field.Length,
			field.Accept,
			field.OnChange,
		)
	}

	m.Form.SetBorder(true).SetTitle("Select Airport")
	m.AddButtons(buttons)

	m.Modal = modal(m.Form, 40, 10)

	return &m
}

func (m *FormModal) SetActionFunc(fn func(int, string)) *FormModal {
	m.ActionFunc = fn
	return m
}

// Adds buttons to the FormModal and sets callbacks for their actions
func (m *FormModal) AddButtons(labels []string) *FormModal {
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

func (m *FormModal) GetInputDataForField(fieldLabel string) string {
	return m.Form.GetFormItemByLabel(fieldLabel).(*tview.InputField).GetText()
}

func (m *FormModal) Primitive() tview.Primitive {
	return m.Modal
}
