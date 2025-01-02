package widgets

import (
	"math"

	"github.com/a-finocchiaro/flightdeck/internal/layout"
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
	// This was calculated manually by figuring out that the max points shown on a graph
	// in the terminal was 171 when tview reported a width of 15.
	//
	// When it's set into a smaller grid component, it's only ~68 and a reported width
	// of 61. This value represents the average of (xWidth * reportedWidth) / 2 to figure
	// out a sensible graph proportion for each screen size.
	pointToWidthProportion = 3356.5
)

type plotData [][]float64

type FlightWidget struct {
	grid     *layout.GridLayout
	Tree     *FlightTree
	gauge    *tvxwidgets.PercentageModeGauge
	altGraph *tvxwidgets.Plot
}

var flightDetailGridOpts layout.GridOptions = layout.GridOptions{
	RowSizes:   []int{0, 3, 0},
	ColSizes:   []int{0},
	HeaderSize: -1,
}

func NewFlightWidget() *FlightWidget {
	g := layout.NewGridLayout(flightDetailGridOpts)

	fw := FlightWidget{
		grid:     g,
		Tree:     NewFlightTree(),
		gauge:    tvxwidgets.NewPercentageModeGauge(),
		altGraph: tvxwidgets.NewPlot(),
	}

	return &fw
}

// Updates the flight view with new data
func (f *FlightWidget) Update(flightData flights.Flight) {
	// clear any existing graph data off of the screen
	f.clearGraphs()

	flightId := flightData.Identification.ID

	// get flight progress information and re-assign flightData if detailed
	// flight information exists.
	if flightId != "" {
		var requester common.Requester = webrequest.SendRequest
		var err error
		flightData, err = client.GetFlightDetails(requester, flightId)

		if err != nil {
			panic(err)
		}

		f.drawFlightProgressBar(flightData)
		f.drawAltitudeGraph(flightData)
	}

	f.Tree.BuildTreeForFlight(flightData)
}

func (f *FlightWidget) Primitive() tview.Primitive {
	f.grid.AddPanel(f.Tree, 0, 0, true)
	f.grid.AddPanel(f.gauge, 1, 0, false)
	f.grid.AddPanel(f.altGraph, 2, 0, false)

	return f.grid
}

// Draws the altitude graph to the screen
func (f *FlightWidget) drawAltitudeGraph(flightData flights.Flight) {
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
func (f *FlightWidget) drawFlightProgressBar(flightData flights.Flight) {
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
func (f *FlightWidget) clearGraphs() {
	f.gauge.SetValue(0)
	f.altGraph.SetData(plotData{})
}
