package widgets

import "github.com/rivo/tview"

// Creates a new Grid and adds the modal content into the grid, and returns the
// resulting Primitive
func NewGridModal(p tview.Primitive, width int, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}
