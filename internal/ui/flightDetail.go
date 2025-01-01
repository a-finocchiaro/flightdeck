package ui

import (
	"github.com/a-finocchiaro/flightdeck/internal/widgets"
)

var FlightDetailPageTitle string = "flightDetail"

type FlightDetailPage struct {
	FlightData *widgets.FlightTree
	Title      string
}

// Constructs a new flight detail page
func NewFlightDetailPage(flightId string) *FlightDetailPage {
	p := FlightDetailPage{
		Title:      FlightDetailPageTitle,
		FlightData: widgets.NewFlightTree(),
	}

	p.FlightData.Update(flightId)

	return &p
}
