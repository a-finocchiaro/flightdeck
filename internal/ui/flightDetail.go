package ui

import (
	"github.com/a-finocchiaro/flightdeck/internal/widgets"
)

var FlightDetailPageTitle string = "flightDetail"

type FlightDetailPage struct {
	FlightData *widgets.FlightWidget
	Title      string
}

// Constructs a new flight detail page
func NewFlightDetailPage() *FlightDetailPage {
	p := FlightDetailPage{
		Title:      FlightDetailPageTitle,
		FlightData: widgets.NewFlightWidget(),
	}

	return &p
}

// Sets the current flight and calls update on the flightdata to fetch an updated
// copy.
func (f *FlightDetailPage) SetFlight(flightId string) {
	f.FlightData.Update(flightId)
}
