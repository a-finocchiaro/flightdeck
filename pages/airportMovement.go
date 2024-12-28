package fr24pages

import (
	"os"

	"github.com/a-finocchiaro/fr24cli/widgets"
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

type AirportMovementPage struct {
	Grid            *tview.Grid
	airport         string
	arrivalTable    *widgets.AirportMovementTable
	departuresTable *widgets.AirportMovementTable
	flightData      *widgets.FlightTree
	airportInfo     *widgets.AirportInfo
	app             *tview.Application
}

// Constructs the new airport movement page
func NewAirportMovementPage(code string, app *tview.Application) *AirportMovementPage {
	page := AirportMovementPage{
		app:             app,
		Grid:            tview.NewGrid(),
		airport:         code,
		arrivalTable:    widgets.NewAirportArrivalsTable(),
		departuresTable: widgets.NewAirportDeparturesTable(),
		flightData:      widgets.NewFlightTree(),
		airportInfo:     widgets.NewAirportInfo(),
	}

	page.buildGrid()
	page.Update()

	page.arrivalTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			app.SetFocus(page.departuresTable)
		}
	})

	page.departuresTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			app.SetFocus(page.arrivalTable)
		}
	})

	page.Grid.AddItem(page.arrivalTable, 1, 0, 1, 1, 0, 100, true).
		AddItem(page.departuresTable, 1, 1, 1, 1, 0, 100, false).
		AddItem(page.flightData.Primitive(), 1, 2, 1, 1, 0, 100, false)

	return &page
}

// fetch updated data and set it in the tables
func (p *AirportMovementPage) Update() {
	// var requester common.Requester = DummyRequester
	var requester common.Requester = webrequest.SendRequest

	airportData, err := client.GetAirportDetails(requester, p.airport, []string{"details"})

	if err != nil {
		panic(err)
	}

	// update the airport info table
	p.airportInfo.Update(airportData)

	// update the arrival and departure tables
	p.arrivalTable.SetData(airportData.Schedule.Arrivals.Data)
	p.departuresTable.SetData(airportData.Schedule.Departures.Data)

	p.arrivalTable.SetSelectedFunc(func(row int, col int) {
		p.flightData.Update(airportData.Schedule.Arrivals.Data[row-1], p.arrivalTable.Table, p.app)
		p.app.SetFocus(p.flightData.Primitive())
	})

	p.departuresTable.SetSelectedFunc(func(row int, col int) {
		p.flightData.Update(airportData.Schedule.Departures.Data[row-1], p.departuresTable.Table, p.app)
		p.app.SetFocus(p.flightData.Primitive())
	})
}

// sets up the base grid
func (p *AirportMovementPage) buildGrid() {
	p.Grid = tview.NewGrid().
		SetRows(10, 0).
		SetColumns(0, 0, 0).
		SetBorders(false).
		AddItem(p.airportInfo.Primitive(), 0, 0, 1, 3, 0, 0, false)
}
