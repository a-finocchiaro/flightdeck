package layout

import "github.com/rivo/tview"

type GridLayout struct {
	*tview.Grid
	opts GridOptions
}

// Sets the sizing options for the grid object. This is expecting an integer array
// where the length is the count of rows or cols, and each value is the size of that
// row or col.
//
// For the HeaderSize value, set this to any size 0 or greater. Setting to -1 will
// disable the header.
type GridOptions struct {
	RowSizes   []int
	ColSizes   []int
	HeaderSize int
}

// Constructor for a new grid layout.
func NewGridLayout(opts GridOptions) *GridLayout {
	rows := opts.RowSizes
	if opts.HeaderSize > -1 {
		rows = append([]int{opts.HeaderSize}, opts.RowSizes...)
	}

	g := GridLayout{
		Grid: tview.NewGrid().SetRows(rows...).SetColumns(opts.ColSizes...),
		opts: opts,
	}

	return &g
}

// Adds a panel to the grid.
func (g *GridLayout) AddPanel(p tview.Primitive, row int, col int, focus bool) {
	g.AddItem(p, row, col, 1, 1, 0, 0, focus)
}

// Adds a pre-built header into the grid and sets the header to span the width of
// the grid.
func (g *GridLayout) AddHeader(p tview.Primitive, focus bool) {
	g.AddItem(p, 0, 0, 1, len(g.opts.ColSizes), 0, 0, focus)
}
