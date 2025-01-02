package ui

import (
	"os"
	"strings"

	"github.com/a-finocchiaro/flightdeck/internal/layout"
	"github.com/a-finocchiaro/flightdeck/internal/widgets"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/client"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/common"
	"github.com/a-finocchiaro/go-flightradar24-sdk/webrequest"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO: remove this
func DummyRequester(s string) ([]byte, error) {
	data, err := os.ReadFile("./tmp/airport_detail.json")

	if err != nil {
		return nil, err
	}

	return data, nil
}

var AirportMovementPageTitle string = "AirportMovementsPage"

type AirportMovementPage struct {
	Grid            *layout.GridLayout
	airport         string
	arrivalTable    *widgets.AirportMovementTable
	departuresTable *widgets.AirportMovementTable
	flightData      *widgets.FlightWidget
	airportInfo     *widgets.AirportInfo
	app             *tview.Application
	pages           *tview.Pages
	Title           string
	Modal           *AirportMovementModal
}

var airportMovementGridOpts layout.GridOptions = layout.GridOptions{
	RowSizes:   []int{0},
	ColSizes:   []int{0, 0, 0},
	HeaderSize: 10,
}

// Constructs the new airport movement page
func NewAirportMovementPage(app *tview.Application, pages *tview.Pages) *AirportMovementPage {
	p := AirportMovementPage{
		Title:           AirportMovementPageTitle,
		app:             app,
		pages:           pages,
		arrivalTable:    widgets.NewAirportArrivalsTable(),
		departuresTable: widgets.NewAirportDeparturesTable(),
		flightData:      widgets.NewFlightWidget(),
		airportInfo:     widgets.NewAirportInfo(),
		Modal:           NewAirportMovementModal(),
	}

	p.Grid = layout.NewGridLayout(airportMovementGridOpts)
	p.Grid.AddHeader(p.airportInfo.Primitive(), false)

	p.arrivalTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			app.SetFocus(p.departuresTable)
		}
	})

	p.departuresTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			app.SetFocus(p.arrivalTable)
		}
	})

	p.Grid.SetBackgroundColor(tcell.ColorDefault)

	p.Grid.AddPanel(p.arrivalTable, 1, 0, true)
	p.Grid.AddPanel(p.departuresTable, 1, 1, true)
	p.Grid.AddPanel(p.flightData.Primitive(), 1, 2, true)

	// set the modal actions
	p.setModalCallback()

	return &p
}

// fetch updated data and set it in the tables
func (p *AirportMovementPage) Update(code string) {
	// var requester common.Requester = DummyRequester
	var requester common.Requester = webrequest.SendRequest

	// update the code
	p.airport = code

	airportData, err := client.GetAirportDetails(
		requester,
		p.airport,
		[]string{"details"},
	)

	if err != nil {
		panic(err)
	}

	// update the airport info table
	p.airportInfo.Update(airportData)

	// update the arrival and departure tables
	p.arrivalTable.SetData(airportData.Schedule.Arrivals.Data)
	p.departuresTable.SetData(airportData.Schedule.Departures.Data)

	p.arrivalTable.SetSelectedFunc(func(row int, col int) {
		data := airportData.Schedule.Arrivals.Data[row-1].Flight.Identification.ID
		p.flightData.Update(data)
		p.setFlightDataEscape(p.arrivalTable.Table)
		p.app.SetFocus(p.flightData.Primitive())
	})

	p.departuresTable.SetSelectedFunc(func(row int, col int) {
		data := airportData.Schedule.Departures.Data[row-1].Flight.Identification.ID
		p.flightData.Update(data)
		p.setFlightDataEscape(p.arrivalTable.Table)
		p.app.SetFocus(p.flightData.Primitive())
	})

	// Send user to the flight detail page if they select a flight
	p.flightData.Tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()

		if ref != nil {
			next := ref.(widgets.NextLocRef)

			// send to the desired airport only if it's a different airport
			if next.Airport != "" && strings.ToLower(next.Airport) != code {
				p.app.SetFocus(p.arrivalTable.Table)
				p.Update(ref.(widgets.NextLocRef).Airport)
			}
		}
	})
}

// sets up the FormModal object to allow input to select an airport
func (p *AirportMovementPage) setModalCallback() {
	p.Modal.SetActionFunc(func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 1 {
			airportCode := p.Modal.GetInputDataForField("Airport IATA:")
			p.Update(airportCode)
			p.pages.ShowPage(p.Title)
			p.pages.HidePage(p.Modal.Title)
		}
		p.pages.HidePage(p.Modal.Title)
	})
}

// Sets the escape to navigate from the the flight data panel back to the tables
func (p *AirportMovementPage) setFlightDataEscape(caller tview.Primitive) {
	p.flightData.Tree.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			p.app.SetFocus(caller)
		}
	})
}
