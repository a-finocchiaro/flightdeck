package utils

import (
	"github.com/a-finocchiaro/flightdeck/internal/styles"
	"github.com/gdamore/tcell/v2"
)

// Determines the status color indicator for a flight based on the Icon value.
func FlightStatusColor(icon string) tcell.Color {
	switch icon {
	case "green":
		return styles.StatusGreen
	case "yellow":
		return styles.StatusYellow
	case "red":
		return styles.StatusRed
	default:
		return styles.StatusUnknown
	}
}
