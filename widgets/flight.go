package widgets

import (
	"fmt"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/airports"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func baseCell(data string) *tview.TableCell {
	return tview.NewTableCell(data).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft)
}

func baseTreeNode(data string) *tview.TreeNode {
	return tview.NewTreeNode(data).SetColor(tcell.ColorWhite)
}

// tree style

type FlightTree struct {
	tree *tview.TreeView
}

func NewFlightTree() *FlightTree {
	fv := FlightTree{
		tree: tview.NewTreeView(),
	}

	return &fv
}

// Updates the flight view with new data
func (f *FlightTree) Update(data airports.FlightArrivalDepartureData, caller tview.Primitive, app *tview.Application) {
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
	baseNode := baseTreeNode(data.Flight.Identification.Callsign)
	f.tree.SetRoot(baseNode).SetCurrentNode(baseNode)
	airportNode := baseTreeNode(direction)
	airlineNode := baseTreeNode("Airline")
	aircraftNode := baseTreeNode("Aircraft")
	statusNode := baseTreeNode("Status")

	baseNode.AddChild(airportNode)
	baseNode.AddChild(airlineNode)
	baseNode.AddChild(aircraftNode)
	baseNode.AddChild(statusNode)

	// Set the Airport info for the origin or destination
	airportNode.AddChild(baseTreeNode(airportData.Name))

	// set the airline info
	airlineNode.AddChild(baseTreeNode(data.Flight.Airline.Name))

	// set the aircraft child node data
	aircraftNode.AddChild(baseTreeNode(data.Flight.Aircraft.Model.Text))
	aircraftNode.AddChild(baseTreeNode(data.Flight.Aircraft.Model.Code))
	aircraftNode.AddChild(baseTreeNode(data.Flight.Aircraft.Registration))

	// update the status indicator
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
			app.SetFocus(caller)
		}
	})
}

func (f *FlightTree) Primitive() tview.Primitive {
	return f.tree
}
