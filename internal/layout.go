package internal

import "github.com/rivo/tview"

type FlightDeckLayout struct {
	app   *tview.Application
	pages *tview.Pages
}

// Constructor for new FlightDeckLayout
func NewLayout(app *tview.Application) *FlightDeckLayout {
	l := FlightDeckLayout{
		app:   app,
		pages: tview.NewPages(),
	}

	return &l
}
