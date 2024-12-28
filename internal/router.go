package internal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Router struct {
	App        *tview.Application
	Pages      *tview.Pages
	Primitives *FlightDeckPrimitives
}

type FlightDeckPrimitives struct {
	AirportMovements *AirportMovementPage
}

// Sets up a new router and starts the application
func Init() {

	pageMgr := FlightDeckPrimitives{}

	r := &Router{
		App:        tview.NewApplication(),
		Pages:      tview.NewPages(),
		Primitives: &pageMgr,
	}

	// add the pages
	r.Primitives.AirportMovements = NewAirportMovementPage(r)

	// add the modal
	r.AddPage("modal", r.Primitives.AirportMovements.Modal().Modal, true, true)

	// bind the keys
	r.bindKeys()

	// start the application
	if err := r.App.SetRoot(r.Pages, true).SetFocus(r.Pages).Run(); err != nil {
		panic(err)
	}
}

// Adds a page to the router
func (r *Router) AddPage(title string, primitive tview.Primitive, resize bool, visible bool) {
	r.Pages.AddPage(title, primitive, resize, visible)
}

// Sets up the application keybindings
func (r *Router) bindKeys() {
	r.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			r.Pages.ShowPage("modal")
			r.Pages.SendToFront("modal")
		}

		return event
	})
}
