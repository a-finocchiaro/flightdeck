package ui

import (
	"github.com/a-finocchiaro/flightdeck/internal/layout"
	"github.com/a-finocchiaro/flightdeck/internal/widgets"
)

var FlightDetailPageTitle string = "flightDetail"

type FlightDetailPage struct {
	Grid       *layout.GridLayout
	flightData *widgets.FlightTree
	title      string
}

var flightDetailGridOpts layout.GridOptions = layout.GridOptions{
	RowCount: 1,
	ColCount: 1,
}

// Constructs a new flight detail page
func NewFlightDetailPage(flightId string) *FlightDetailPage {
	p := FlightDetailPage{
		title:      FlightDetailPageTitle,
		flightData: widgets.NewFlightTree(),
	}

	p.Grid = layout.NewGridLayout(flightDetailGridOpts)
	p.Grid.AddPanel(p.flightData.Primitive(), 0, 0, true)
	p.flightData.Update(flightId)

	return &p
}
