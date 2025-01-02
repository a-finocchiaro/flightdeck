package widgets

import (
	"fmt"

	"github.com/a-finocchiaro/go-flightradar24-sdk/pkg/models/flights"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FlightTree struct {
	*tview.TreeView
	RootNode     *tview.TreeNode
	StatusNode   *tview.TreeNode
	AirlineNode  *tview.TreeNode
	OriginNode   *tview.TreeNode
	DestNode     *tview.TreeNode
	AircraftNode *tview.TreeNode
	SpeedAltNode *tview.TreeNode
}

type NextLocRef struct {
	Airport string
}

// Builds a base tree node
func baseTreeNode(data string) *tview.TreeNode {
	return tview.NewTreeNode(data).SetColor(tcell.ColorWhite)
}

// Assembles a new flight tree
func NewFlightTree() *FlightTree {
	ft := FlightTree{
		TreeView: tview.NewTreeView(),
	}

	return &ft
}

// constructs a flight tree from flight data
func (t *FlightTree) BuildTreeForFlight(flight flights.Flight) {
	t.RootNode = baseTreeNode(flight.Identification.Number.Default)
	t.TreeView.SetRoot(t.RootNode).SetCurrentNode(t.RootNode)
	t.AirlineNode = baseTreeNode("Airline")
	t.OriginNode = baseTreeNode("Origin")
	t.DestNode = baseTreeNode("Destination")
	t.AircraftNode = baseTreeNode("Aircraft")
	t.StatusNode = baseTreeNode("Status")
	t.SpeedAltNode = baseTreeNode("Speed and Altitude")

	t.RootNode.AddChild(t.StatusNode)
	t.RootNode.AddChild(t.OriginNode)
	t.RootNode.AddChild(t.DestNode)
	t.RootNode.AddChild(t.AirlineNode)
	t.RootNode.AddChild(t.AircraftNode)
	t.RootNode.AddChild(t.SpeedAltNode)

	// Set the Airport info for the origin or destination
	origin := baseTreeNode(flight.Airport.Origin.Name)
	origin.SetReference(NextLocRef{Airport: flight.Airport.Origin.Code.Iata})
	t.OriginNode.AddChild(origin)

	dest := baseTreeNode(flight.Airport.Destination.Name)
	dest.SetReference(NextLocRef{Airport: flight.Airport.Destination.Code.Iata})
	t.DestNode.AddChild(dest)

	// set the airline info
	t.AirlineNode.AddChild(baseTreeNode(flight.Airline.Name))

	// set the aircraft child node data
	t.AircraftNode.AddChild(baseTreeNode(flight.Aircraft.Model.Text))
	t.AircraftNode.AddChild(baseTreeNode(flight.Aircraft.Model.Code))
	t.AircraftNode.AddChild(baseTreeNode(flight.Aircraft.Registration))

	var statusColor tcell.Color

	switch flight.Status.Icon {
	case "green":
		statusColor = tcell.ColorGreen
	case "yellow":
		statusColor = tcell.ColorYellow
	case "red":
		statusColor = tcell.ColorRed
	default:
		statusColor = tcell.ColorGray
	}

	t.StatusNode.AddChild(baseTreeNode(fmt.Sprintf("%s %s", "⏺", flight.Status.Text)).SetColor(statusColor))

	// speed and altitude data
	if len(flight.Trail) > 0 {
		latestTrail := flight.Trail[0]
		t.SpeedAltNode.AddChild(baseTreeNode(fmt.Sprintf("%dft", latestTrail.Alt)))
		t.SpeedAltNode.AddChild(baseTreeNode(fmt.Sprintf("%dkts", latestTrail.Spd)))
	}
}
