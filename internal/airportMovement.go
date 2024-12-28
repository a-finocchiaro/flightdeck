package internal

import (
	"os"

	"github.com/a-finocchiaro/flightdeck/widgets"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/client"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/common"
	"github.com/gdamore/tcell/v2"
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
	Grid            *GridLayout
	airport         string
	arrivalTable    *widgets.AirportMovementTable
	departuresTable *widgets.AirportMovementTable
	flightData      *widgets.FlightTree
	airportInfo     *widgets.AirportInfo
	router          *Router
	title           string
	Modal           *AirportMovementModal
}

var airportMovementGridOpts GridOptions = GridOptions{
	RowCount:   1,
	ColCount:   3,
	HeaderSize: 10,
}

// Constructs the new airport movement page
func NewAirportMovementPage(router *Router) *AirportMovementPage {
	p := AirportMovementPage{
		title:           AirportMovementPageTitle,
		router:          router,
		arrivalTable:    widgets.NewAirportArrivalsTable(),
		departuresTable: widgets.NewAirportDeparturesTable(),
		flightData:      widgets.NewFlightTree(router.App),
		airportInfo:     widgets.NewAirportInfo(),
		Modal:           NewAirportMovementModal(),
	}

	p.Grid = NewGridLayout(airportMovementGridOpts)
	p.Grid.AddHeader(p.airportInfo.Primitive(), false)

	p.arrivalTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			router.App.SetFocus(p.departuresTable)
		}
	})

	p.departuresTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			router.App.SetFocus(p.arrivalTable)
		}
	})

	p.Grid.AddPanel(p.arrivalTable, 1, 0, true)
	p.Grid.AddPanel(p.departuresTable, 1, 1, true)
	p.Grid.AddPanel(p.flightData.Primitive(), 1, 2, true)

	// set the modal actions
	p.setModalCallback()

	router.AddPage(p.Modal.Title, p.Modal.Primitive(), true, true)
	router.AddPage(p.title, p.Grid, true, false)

	return &p
}

// fetch updated data and set it in the tables
func (p *AirportMovementPage) Update(code string) {
	var requester common.Requester = DummyRequester
	// var requester common.Requester = webrequest.SendRequest

	// update the code
	p.airport = code

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
		p.flightData.Update(airportData.Schedule.Arrivals.Data[row-1], p.arrivalTable.Table)
		p.router.App.SetFocus(p.flightData.Primitive())
	})

	p.departuresTable.SetSelectedFunc(func(row int, col int) {
		p.flightData.Update(airportData.Schedule.Departures.Data[row-1], p.departuresTable.Table)
		p.router.App.SetFocus(p.flightData.Primitive())
	})
}

// sets up the FormModal object to allow input to select an airport
func (p *AirportMovementPage) setModalCallback() {
	p.Modal.SetActionFunc(func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 1 {
			airportCode := p.Modal.GetInputDataForField("Airport IATA:")
			p.Update(airportCode)
			p.router.Pages.ShowPage(p.title)
			p.router.Pages.HidePage(p.Modal.Title)
		}
		p.router.Pages.HidePage(p.Modal.Title)
	})
}
