package config

const (
	DefaultAirport = ""
)

// Type associated with possible configuration values in FlightDeck
type FlightDeckConfig struct {
	Airport string
}

// instantiates a new flight deck config object and populates with default vals
func NewFlightDeckConfig() *FlightDeckConfig {
	return &FlightDeckConfig{
		Airport: DefaultAirport,
	}
}
