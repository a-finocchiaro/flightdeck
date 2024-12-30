package widgets

import (
	"fmt"
	"math"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/client"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/airports"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/common"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/flights"
	"github.com/a-finocchiaro/go-flightradar24-sdk/webrequest"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"github.com/umahmood/haversine"
)

type plotData [][]float64

func baseTreeNode(data string) *tview.TreeNode {
	return tview.NewTreeNode(data).SetColor(tcell.ColorWhite)
}

type FlightTree struct {
	grid     *tview.Grid
	tree     *tview.TreeView
	gauge    *tvxwidgets.PercentageModeGauge
	altGraph *tvxwidgets.Plot
	app      *tview.Application
}

func NewFlightTree(app *tview.Application) *FlightTree {
	g := tview.NewGrid().SetRows(0, 3, 0).SetColumns(0)

	fv := FlightTree{
		grid:     g,
		tree:     tview.NewTreeView(),
		gauge:    tvxwidgets.NewPercentageModeGauge(),
		altGraph: tvxwidgets.NewPlot(),
		app:      app,
	}

	return &fv
}

// Updates the flight view with new data
func (f *FlightTree) Update(data airports.FlightArrivalDepartureData, caller tview.Primitive) {
	// clear any existing graph data off of the screen
	f.clearGraphs()

	// set the airport data since this is stored either under arrival or departure
	var airportData airports.FlightAiportData
	var direction string

	if data.Flight.Status.Generic.Status.Type == "arrival" {
		airportData = data.Flight.Airport.Origin
		direction = "Origin"
	} else {
		airportData = data.Flight.Airport.Destination
		direction = "Destination"
	}

	// Setup the base tree
	baseNode := baseTreeNode(data.Flight.Identification.Number.Default)
	f.tree.SetRoot(baseNode).SetCurrentNode(baseNode)
	airportNode := baseTreeNode(direction)
	airlineNode := baseTreeNode("Airline")
	aircraftNode := baseTreeNode("Aircraft")
	statusNode := baseTreeNode("Status")

	baseNode.AddChild(statusNode)
	baseNode.AddChild(airportNode)
	baseNode.AddChild(airlineNode)
	baseNode.AddChild(aircraftNode)

	// Set the Airport info for the origin or destination
	airportNode.AddChild(baseTreeNode(airportData.Name))

	// set the airline info
	airlineNode.AddChild(baseTreeNode(data.Flight.Airline.Name))

	// set the aircraft child node data
	aircraftNode.AddChild(baseTreeNode(data.Flight.Aircraft.Model.Text))
	aircraftNode.AddChild(baseTreeNode(data.Flight.Aircraft.Model.Code))
	aircraftNode.AddChild(baseTreeNode(data.Flight.Aircraft.Registration))

	// set the status indicator
	var statusColor tcell.Color

	switch data.Flight.Status.Icon {
	case "green":
		statusColor = tcell.ColorGreen
	case "yellow":
		statusColor = tcell.ColorYellow
	case "red":
		statusColor = tcell.ColorRed
	default:
		statusColor = tcell.ColorGray
	}

	statusNode.AddChild(baseTreeNode(fmt.Sprintf("%s %s", "⏺", data.Flight.Status.Text)).SetColor(statusColor))

	// set the escape
	f.tree.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			f.app.SetFocus(caller)
		}
	})

	// get flight progress information
	var requester common.Requester = webrequest.SendRequest
	flightData, err := client.GetFlightDetails(requester, data.Flight.Identification.ID)

	if err != nil {
		panic(err)
	}

	f.drawFlightProgressBar(flightData)
	f.drawAltitudeGraph(flightData)
}

func (f *FlightTree) Primitive() tview.Primitive {
	f.grid.AddItem(f.tree, 0, 0, 1, 1, 0, 0, true).
		AddItem(f.gauge, 1, 0, 1, 1, 0, 0, false).
		AddItem(f.altGraph, 2, 0, 1, 1, 0, 0, false)

	return f.grid
}

// Draws the altitude graph to the screen
func (f *FlightTree) drawAltitudeGraph(flightData flights.Flight) {
	trail := flightData.Trail

	// Do nothing for flights that have not departed yet or do not contain any
	// trail data.
	if len(trail) == 0 {
		return
	}

	altData := plotData{{0}}

	// find the sample rate based on the width of the terminal. If your
	// terminal window is tiny, then yes the graph will be kind of useless
	_, _, termWidth, _ := f.altGraph.Box.GetInnerRect()

	q := float64(len(trail)) / float64(termWidth)
	sampleRate := int(math.Ceil(q))

	// insert the altitude records into the data at the sample rate that matches
	for i := len(trail) - 1; i >= 0; i-- {

		if i%sampleRate != 0 {
			continue
		}

		alt := trail[i].Alt
		altData[0] = append(altData[0], float64(alt))
	}

	f.altGraph.SetMarker(tvxwidgets.PlotMarkerBraille)
	f.altGraph.SetDrawXAxisLabel(true)
	f.altGraph.SetTitle("Altitude Graph")
	f.altGraph.SetBorder(true)
	f.altGraph.SetLineColor([]tcell.Color{tcell.ColorSteelBlue})
	f.altGraph.SetData(altData)
}

// Draws the flight progress bar to the screen
func (f *FlightTree) drawFlightProgressBar(flightData flights.Flight) {
	trail := flightData.Trail

	// Do nothing for flights that have not departed yet or do not contain any
	// trail data.
	if len(trail) == 0 {
		return
	}

	currCoord := haversine.Coord{
		Lat: trail[0].Lat,
		Lon: trail[0].Lng,
	}

	originCoord := haversine.Coord{
		Lat: flightData.Airport.Origin.Position.Latitude,
		Lon: flightData.Airport.Origin.Position.Longitude,
	}

	destCoord := haversine.Coord{
		Lat: flightData.Airport.Destination.Position.Latitude,
		Lon: flightData.Airport.Destination.Position.Longitude,
	}

	totalDistMi, _ := haversine.Distance(originCoord, destCoord)
	remainingDist, _ := haversine.Distance(currCoord, destCoord)

	f.gauge.SetTitle("Flight Progress")
	f.gauge.SetBorder(true)
	f.gauge.SetMaxValue(int(totalDistMi))
	f.gauge.SetValue(int(totalDistMi) - int(remainingDist))

	f.drawAltitudeGraph(flightData)
}

// Automatically clears the data out of all graphs
func (f *FlightTree) clearGraphs() {
	f.gauge.SetValue(0)
	f.altGraph.SetData(plotData{})
}
