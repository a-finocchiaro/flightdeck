package ui

import "github.com/a-finocchiaro/flightdeck/internal/widgets"

var AirportMovementModalPageTitle string = "AirportMovementsModal"

type AirportMovementModal struct {
	*widgets.FormModal
	Title string
}

func NewAirportMovementModal() *AirportMovementModal {
	// setup the input fields
	iataCode := []widgets.InputFields{
		{
			Label:       "Airport IATA:",
			Placeholder: "",
			Length:      5,
			Accept:      nil,
			OnChange:    nil,
		},
	}

	buttons := []string{"Cancel", "Accept"}

	m := &AirportMovementModal{
		FormModal: widgets.NewFormModal(buttons, iataCode),
		Title:     AirportMovementModalPageTitle,
	}

	return m
}
