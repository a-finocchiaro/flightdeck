package widgets

import (
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/airports"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AirportInfo struct {
	table *tview.Table
}

func baseCell(data string) *tview.TableCell {
	return tview.NewTableCell(data).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft)
}

// Constructs a new AirportInfo object
func NewAirportInfo() *AirportInfo {
	airportInfo := AirportInfo{
		table: tview.NewTable().SetSelectable(false, false),
	}
	airportInfo.baseLayout()

	return &airportInfo
}

// Updates the flight view with new data
func (a *AirportInfo) Update(data airports.AirportPluginData) {
	a.table.SetTitle(data.Details.Name)
	a.table.SetBorder(true)
	a.table.SetCell(1, 1, baseCell(data.Details.Code.Iata))
	a.table.SetCell(2, 1, baseCell(data.Details.Code.Icao))
	a.table.SetCell(3, 1, baseCell(data.Weather.Metar))
}

// Returns a tview Primitive to display
func (a *AirportInfo) Primitive() tview.Primitive {
	return a.table
}

// Creates the base table layout
func (a *AirportInfo) baseLayout() {
	a.table.SetCell(1, 0, baseCell("IATA:"))
	a.table.SetCell(2, 0, baseCell("ICAO:"))
	a.table.SetCell(3, 0, baseCell("Current Metar: "))
	a.table.SetCell(1, 2, baseCell("Arrival Delay Index: "))
	a.table.SetCell(2, 2, baseCell("Departure Delay Index: "))
	a.table.SetCell(1, 3, baseCell("Arrival Delay Avg: "))
	a.table.SetCell(2, 3, baseCell("Departure Delay Avg: "))
}
