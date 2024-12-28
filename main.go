package main

import (
	fr24pages "github.com/a-finocchiaro/fr24cli/pages"
	"github.com/a-finocchiaro/fr24cli/widgets"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	modal := widgets.NewFr24Modal("What do we set this to?")

	modal.SetActionFunc(func(buttonIndex int, buttonLabel string) {
		// set the page to be the new airport
		if buttonIndex == 1 {
			airportCode := modal.Form.GetFormItemByLabel("Airport IATA:").(*tview.InputField).GetText()
			airportMovement := fr24pages.NewAirportMovementPage(airportCode, app)
			pages.AddPage("Airport Movements", airportMovement.Grid, true, true)
			pages.HidePage("modal")
		}
		pages.HidePage("modal")
	})

	pages.AddPage("modal", modal.Modal, true, true)

	// // set the input capture
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			pages.ShowPage("modal")
			pages.SendToFront("modal")
		}

		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}

}
