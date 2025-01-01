package internal

import (
	"strings"

	"github.com/a-finocchiaro/flightdeck/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Router struct {
	App        *tview.Application
	Pages      *tview.Pages
	Primitives *FlightDeckPrimitives
}

type FlightDeckPrimitives struct {
	AirportMovements *ui.AirportMovementPage
	HelpModal        *HelpModal
	FlightDetailPage *FlightDetailPage
}

const (
	KeyHelp tcell.Key = 63
)

func initKeys() {
	tcell.KeyNames[KeyHelp] = "?"
}

// Sets up a new router and starts the application
func Init() {
	initKeys()
	pageMgr := FlightDeckPrimitives{}

	r := &Router{
		App:        tview.NewApplication(),
		Pages:      tview.NewPages(),
		Primitives: &pageMgr,
	}

	// add the pages
	r.Primitives.AirportMovements = ui.NewAirportMovementPage(r.App, r.Pages)
	r.Primitives.HelpModal = NewHelpModal(r)
	// r.Primitives.FlightDetailPage = NewFlightDetailPage(r, "388bb127")
	r.AddPage(
		r.Primitives.AirportMovements.Modal.Title,
		r.Primitives.AirportMovements.Modal.Primitive(),
		true,
		true,
	)
	r.AddPage(
		r.Primitives.AirportMovements.Title,
		r.Primitives.AirportMovements.Grid,
		true,
		false,
	)

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
		switch event.Key() {

		// Pressing F1 opens the Airport select modal
		case tcell.KeyF1:
			r.Pages.ShowPage(r.Primitives.AirportMovements.Modal.Title)
			r.Pages.SendToFront(r.Primitives.AirportMovements.Modal.Title)

		// escape should always close a modal if one is open
		case tcell.KeyEsc:
			name, _ := r.Pages.GetFrontPage()

			if strings.Contains(name, "Modal") {
				r.Pages.HidePage(name)
			}

		case tcell.KeyRune:
			if event.Rune() == '?' {
				r.Pages.ShowPage(r.Primitives.HelpModal.Title)
				r.Pages.SendToFront(r.Primitives.HelpModal.Title)
			}
		}

		return event
	})
}
