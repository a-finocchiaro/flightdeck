package internal

import (
	"github.com/a-finocchiaro/flightdeck/widgets"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var HelpModalPageName = "HelpModal"

type HelpModal struct {
	Modal *widgets.StaticModal
	Title string
	table *tview.Table
}

type HelpHint struct {
	Key         string
	Description string
}

type HelpHints []HelpHint

// Creates a new help modal
func NewHelpModal(r *Router) *HelpModal {
	t := tview.NewTable()

	// setup a table with some help information
	for idx, hint := range buildHelpHints() {
		keyCellVal := tview.NewTableCell(hint.Key).SetAlign(tview.AlignLeft).SetTextColor(tcell.ColorYellow)
		expCellVal := tview.NewTableCell(hint.Description).SetAlign(tview.AlignLeft)

		t.SetCell(idx, 0, keyCellVal)
		t.SetCell(idx, 1, expCellVal)
	}

	t.SetBorder(true).SetTitle("Help")

	m := HelpModal{
		table: t,
		Title: HelpModalPageName,
	}

	m.Modal = widgets.NewStaticModal(t)
	r.AddPage(HelpModalPageName, m.Modal.Primitive(), true, false)

	return &m
}

// Builds a collection of HelpHint objects. Whenever a new keyboard shortcut is added,
// an explanation of that shortcut should be added here.
func buildHelpHints() HelpHints {
	return HelpHints{
		{
			Key:         "F1",
			Description: "Select Airport",
		},
		{
			Key:         "Esc",
			Description: "Close Popup Window",
		},
		{
			Key:         "Ctrl+C",
			Description: "Exit Application",
		},
		{
			Key:         "?",
			Description: "Open this help screen",
		},
	}
}
