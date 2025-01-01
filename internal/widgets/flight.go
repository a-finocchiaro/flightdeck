package widgets

import (
	"fmt"
	"math"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/client"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/common"
	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/flights"
	"github.com/a-finocchiaro/go-flightradar24-sdk/webrequest"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"github.com/umahmood/haversine"
)

const (
	// This was calculated manually by figuring out that the max points on a graph at
	// full screen is 255 with a reported width of 15 from tview.
	//
	// When it's set into a smaller grid component, it's only ~68 and a reported width
	// of 61. This value represents the average of (xWidth x reportedWidth) to figure
	// out a sensible graph proportion for each screen size.
	pointToWidthProportion = 3986.5
)

type plotData [][]float64

func baseTreeNode(data string) *tview.TreeNode {
	return tview.NewTreeNode(data).SetColor(tcell.ColorWhite)
}

type FlightTree struct {
	grid     *tview.Grid
	Tree     *tview.TreeView
	gauge    *tvxwidgets.PercentageModeGauge
	altGraph *tvxwidgets.Plot
}

func NewFlightTree() *FlightTree {
	g := tview.NewGrid().SetRows(0, 3, 0).SetColumns(0)

	fv := FlightTree{
		grid:     g,
		Tree:     tview.NewTreeView(),
		gauge:    tvxwidgets.NewPercentageModeGauge(),
		altGraph: tvxwidgets.NewPlot(),
	}

	return &fv
}

// Updates the flight view with new data
func (f *FlightTree) Update(flightId string) {
	// clear any existing graph data off of the screen
	f.clearGraphs()

	// get flight progress information
	var requester common.Requester = webrequest.SendRequest
	flightData, err := client.GetFlightDetails(requester, flightId)

	if err != nil {
		panic(err)
	}

	// Setup the base tree
	baseNode := baseTreeNode(flightData.Identification.Number.Default)
	f.Tree.SetRoot(baseNode).SetCurrentNode(baseNode)
	airlineNode := baseTreeNode("Airline")
	originNode := baseTreeNode("Origin")
	destNode := baseTreeNode("Destination")
	aircraftNode := baseTreeNode("Aircraft")
	statusNode := baseTreeNode("Status")

	baseNode.AddChild(statusNode)
	baseNode.AddChild(originNode)
	baseNode.AddChild(destNode)
	baseNode.AddChild(airlineNode)
	baseNode.AddChild(aircraftNode)

	// Set the Airport info for the origin or destination
	originNode.AddChild(baseTreeNode(flightData.Airport.Origin.Name))
	destNode.AddChild(baseTreeNode(flightData.Airport.Destination.Name))

	// set the airline info
	airlineNode.AddChild(baseTreeNode(flightData.Airline.Name))

	// set the aircraft child node data
	aircraftNode.AddChild(baseTreeNode(flightData.Aircraft.Model.Text))
	aircraftNode.AddChild(baseTreeNode(flightData.Aircraft.Model.Code))
	aircraftNode.AddChild(baseTreeNode(flightData.Aircraft.Registration))

	// set the status indicator
	var statusColor tcell.Color

	switch flightData.Status.Icon {
	case "green":
		statusColor = tcell.ColorGreen
	case "yellow":
		statusColor = tcell.ColorYellow
	case "red":
		statusColor = tcell.ColorRed
	default:
		statusColor = tcell.ColorGray
	}

	statusNode.AddChild(baseTreeNode(fmt.Sprintf("%s %s", "⏺", flightData.Status.Text)).SetColor(statusColor))

	f.drawFlightProgressBar(flightData)
	f.drawAltitudeGraph(flightData)
}

func (f *FlightTree) Primitive() tview.Primitive {
	f.grid.AddItem(f.Tree, 0, 0, 1, 1, 0, 0, true).
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
	pval := pointToWidthProportion / float64(termWidth)

	q := float64(len(trail)) / float64(pval)
	sampleRate := int(math.Ceil(q))
	// sampleRate := termWidth

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