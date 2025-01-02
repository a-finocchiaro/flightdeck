package widgets

import (
	"errors"
	"strings"
	"time"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/airports"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/flights"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var tableHeaders = []string{"Flight No", "Airline", "Time"}

// Integer alias to indicate the type of movement table.
type MovementState int

// constants with iota values used to indicate arrivals (0) or departures(1).
const (
	arrivals MovementState = iota
	departures
)

type AirportMovementTable struct {
	*tview.Table
	movementState MovementState
}

func NewAirportArrivalsTable() *AirportMovementTable {
	table, _ := New("arrivals")

	return table
}

func NewAirportDeparturesTable() *AirportMovementTable {
	table, _ := New("departures")

	return table
}

func New(title string) (*AirportMovementTable, error) {
	// validate the title to make sure it is either arrivals or departures, these are
	// the only movement types we allow.
	var movementState MovementState

	switch strings.ToLower(title) {
	case "arrivals":
		movementState = arrivals
	case "departures":
		movementState = departures
	default:
		return nil, errors.New("title must be arrivals or departures")
	}

	t := &AirportMovementTable{
		Table:         tview.NewTable(),
		movementState: movementState,
	}

	t.baseLayout()
	t.SetBorders(false).SetSelectable(true, false).Select(1, 0).SetFixed(1, 1).SetBorder(true)
	t.SetTitle(title)
	return t, nil
}

// Sets the movement (arrival or departure) data into the table
func (t *AirportMovementTable) SetData(data []airports.FlightArrivalDepartureData) {
	for row, flight := range data {
		var returnCity flights.FlightAirport

		// set the return city value
		if t.movementState == arrivals {
			returnCity = flight.Flight.Airport.Origin
		} else {
			returnCity = flight.Flight.Airport.Destination
		}

		t.SetCell(row+1, 0, tview.NewTableCell(returnCity.Code.Iata).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
		t.SetCell(row+1, 1, tview.NewTableCell(flight.Flight.Identification.Number.Default).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))
		t.SetCell(row+1, 2, tview.NewTableCell(flight.Flight.Airline.Short).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter))

		// convert the time to a string
		strTime := time.Unix(int64(flight.Flight.Time.Scheduled.Arrival), 0).String()

		t.SetCell(row+1, 3, tview.NewTableCell(strTime).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter).SetExpansion(1))
	}
}

// Sets the table headers depending on the type of movement table that is being built.
func (t *AirportMovementTable) baseLayout() {
	var headers []string

	// Set the first column either as Origin or Destination depending on the movement
	// type of the table.
	if t.movementState == arrivals {
		headers = append([]string{"Origin"}, tableHeaders...)
	} else {
		headers = append([]string{"Destination"}, tableHeaders...)
	}

	for i, h := range headers {
		cellVal := tview.NewTableCell(h).SetTextColor(tcell.ColorWhite).SetSelectable(false).SetAlign(tview.AlignCenter)
		t.SetCell(0, i, cellVal)
	}
}
