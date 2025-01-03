package widgets

import (
	"fmt"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/airports"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AirportInfo struct {
	grid                tview.Grid
	infoTable           *tview.Table
	weatherTable        *tview.Table
	arrivalStatsTable   *tview.Table
	departureStatsTable *tview.Table
}

func baseCell(data string) *tview.TableCell {
	return tview.NewTableCell(data).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft)
}

// Constructs a new AirportInfo object
func NewAirportInfo() *AirportInfo {
	airportInfo := AirportInfo{
		grid:                *tview.NewGrid().SetRows([]int{0}...).SetColumns([]int{0, 0, 0, 0}...),
		infoTable:           tview.NewTable(),
		weatherTable:        tview.NewTable(),
		arrivalStatsTable:   tview.NewTable(),
		departureStatsTable: tview.NewTable(),
	}
	airportInfo.baseLayout()
	airportInfo.grid.SetBorder(true)
	airportInfo.arrivalStatsTable.SetBorder(true).SetTitle("Arrival Data")
	airportInfo.departureStatsTable.SetBorder(true).SetTitle("Departure Data")
	airportInfo.grid.AddItem(airportInfo.infoTable, 0, 0, 1, 1, 0, 0, false)
	airportInfo.grid.AddItem(airportInfo.weatherTable, 0, 1, 1, 1, 0, 0, false)
	airportInfo.grid.AddItem(airportInfo.arrivalStatsTable, 0, 2, 1, 1, 0, 0, false)
	airportInfo.grid.AddItem(airportInfo.departureStatsTable, 0, 3, 1, 1, 0, 0, false)

	return &airportInfo
}

// Updates the flight view with new data
func (a *AirportInfo) Update(data airports.AirportPluginData) {
	a.grid.SetTitle(data.Details.Name)
	// a.infoTable.SetBorder(true)
	a.infoTable.SetCell(1, 1, baseCell(data.Details.Code.Iata))
	a.infoTable.SetCell(2, 1, baseCell(data.Details.Code.Icao))
	a.infoTable.SetCell(3, 1, baseCell(fmt.Sprintf("%d", data.Weather.Elevation.Ft)))

	// weather
	a.weatherTable.SetCell(1, 1, baseCell(fmt.Sprintf("%d", data.Weather.Temp.Fahrenheit)))
	a.weatherTable.SetCell(2, 1, baseCell(fmt.Sprintf("%d", data.Weather.Humidity)))
	a.weatherTable.SetCell(3, 1, baseCell(data.Weather.Sky.Condition.Text))
	a.weatherTable.SetCell(4, 1, baseCell(fmt.Sprintf("%d", data.Weather.Sky.Visibility.Nmi)))
	a.weatherTable.SetCell(5, 1, baseCell(fmt.Sprintf("%d@%d", data.Weather.Wind.Direction.Degree, data.Weather.Wind.Speed.Kts)))

	// Arrivals Stats
	a.arrivalStatsTable.SetCell(1, 1, baseCell(fmt.Sprintf("%.2f", data.Details.Stats.Arrivals.DelayIndex)))
	a.arrivalStatsTable.SetCell(2, 1, baseCell(fmt.Sprintf("%d", data.Details.Stats.Arrivals.DelayAvg)))
	a.arrivalStatsTable.SetCell(3, 1, baseCell(fmt.Sprintf("%d", data.Details.Stats.Arrivals.Today.Quantity.OnTime)))
	a.arrivalStatsTable.SetCell(4, 1, baseCell(fmt.Sprintf("%d", data.Details.Stats.Arrivals.Today.Quantity.Delayed)))
	a.arrivalStatsTable.SetCell(5, 1, baseCell(fmt.Sprintf("%d", data.Details.Stats.Arrivals.Today.Quantity.Canceled)))

	// Departure Stats
	a.departureStatsTable.SetCell(1, 1, baseCell(fmt.Sprintf("%.2f", data.Details.Stats.Departures.DelayIndex)))
	a.departureStatsTable.SetCell(2, 1, baseCell(fmt.Sprintf("%d", data.Details.Stats.Departures.DelayAvg)))
	a.departureStatsTable.SetCell(3, 1, baseCell(fmt.Sprintf("%d", data.Details.Stats.Departures.Today.Quantity.OnTime)))
	a.departureStatsTable.SetCell(4, 1, baseCell(fmt.Sprintf("%d", data.Details.Stats.Departures.Today.Quantity.Delayed)))
	a.departureStatsTable.SetCell(5, 1, baseCell(fmt.Sprintf("%d", data.Details.Stats.Departures.Today.Quantity.Canceled)))
}

// Returns a tview Primitive to display
func (a *AirportInfo) Primitive() tview.Primitive {
	return &a.grid
}

// Creates the base table layout
func (a *AirportInfo) baseLayout() {
	a.infoTable.SetCell(1, 0, baseCell("IATA:"))
	a.infoTable.SetCell(2, 0, baseCell("ICAO:"))
	a.infoTable.SetCell(3, 0, baseCell("Elevation (ft):"))
	a.infoTable.SetCell(4, 0, baseCell(""))
	a.infoTable.SetCell(5, 0, baseCell(""))

	// weather col
	a.weatherTable.SetCell(1, 0, baseCell("Temp (F):"))
	a.weatherTable.SetCell(2, 0, baseCell("Humidity:"))
	a.weatherTable.SetCell(3, 0, baseCell("Conditions:"))
	a.weatherTable.SetCell(4, 0, baseCell("Visibility (nmi):"))
	a.weatherTable.SetCell(5, 0, baseCell("Wind (kts):"))

	// arrival stats
	a.arrivalStatsTable.SetCell(1, 0, baseCell("Delay Index:"))
	a.arrivalStatsTable.SetCell(2, 0, baseCell("Delay Avg:"))
	a.arrivalStatsTable.SetCell(3, 0, baseCell("Total On Time Today:"))
	a.arrivalStatsTable.SetCell(4, 0, baseCell("Total Delayed Today:"))
	a.arrivalStatsTable.SetCell(5, 0, baseCell("Total Canceled Today:"))

	// departure stats
	a.departureStatsTable.SetCell(1, 0, baseCell("Delay Index:"))
	a.departureStatsTable.SetCell(2, 0, baseCell("Delay Avg:"))
	a.departureStatsTable.SetCell(3, 0, baseCell("Total On Time Today:"))
	a.departureStatsTable.SetCell(4, 0, baseCell("Total Delayed Today:"))
	a.departureStatsTable.SetCell(5, 0, baseCell("Total Canceled Today:"))
}
