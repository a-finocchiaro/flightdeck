package internal

import (
	"strings"
	"time"

	"github.com/a-finocchiaro/flightdeck/config"
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
	HelpModal        *ui.HelpModal
}

// Sets up a new router and starts the application
func Init(cfg *config.FlightDeckConfig) {
	pageMgr := FlightDeckPrimitives{}

	r := &Router{
		App:        tview.NewApplication(),
		Pages:      tview.NewPages(),
		Primitives: &pageMgr,
	}

	// add the pages
	r.Primitives.AirportMovements = ui.NewAirportMovementPage(r.App, r.Pages)
	r.Primitives.HelpModal = ui.NewHelpModal()

	r.AddPage(r.Primitives.HelpModal.Title, r.Primitives.HelpModal.Modal.Primitive(), true, true)
	r.AddPage(
		r.Primitives.AirportMovements.Modal.Title,
		r.Primitives.AirportMovements.Modal.Primitive(),
		true,
		false,
	)
	r.AddPage(
		r.Primitives.AirportMovements.Title,
		r.Primitives.AirportMovements.Grid,
		true,
		false,
	)

	// if an airport is defined from config, open directly to that airport
	// TODO: workflow here is a little jank, maybe move/refactor this?
	// Tip: it's jank because we set the visibility of the windows above, then re-hide
	// them here.
	if cfg.Airport != "" {
		r.Pages.HidePage(r.Primitives.HelpModal.Title)
		r.Primitives.AirportMovements.Start(cfg.Airport)
		r.Pages.ShowPage(r.Primitives.AirportMovements.Title)
	}

	// bind the keys
	r.bindKeys()

	// Add the automatic refresh loop
	// TODO: make this configurable, and maybe add ability to disable?
	go func() {
		for {
			time.Sleep(5 * time.Second)
			r.Primitives.AirportMovements.Update()
			r.App.Draw()
		}
	}()

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
