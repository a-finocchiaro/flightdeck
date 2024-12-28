package internal

import "github.com/rivo/tview"

type GridLayout struct {
	*tview.Grid
	opts GridOptions
}

type GridOptions struct {
	RowCount   int
	ColCount   int
	HeaderSize int
}

func NewGridLayout(opts GridOptions) *GridLayout {
	rows := make([]int, opts.RowCount)
	cols := make([]int, opts.ColCount)

	if opts.HeaderSize > -1 {
		rows = append([]int{opts.HeaderSize}, rows...)
	}

	g := GridLayout{
		Grid: tview.NewGrid().SetRows(rows...).SetColumns(cols...),
		opts: opts,
	}

	return &g
}

func (g *GridLayout) AddPanel(p tview.Primitive, row int, col int, focus bool) {
	g.AddItem(p, row, col, 1, 1, 0, 100, focus)
}

func (g *GridLayout) AddHeader(p tview.Primitive, focus bool) {
	g.AddItem(p, 0, 0, 1, g.opts.ColCount, 0, 0, focus)
}
