package cmd

import (
	"github.com/a-finocchiaro/flightdeck/config"
	"github.com/a-finocchiaro/flightdeck/internal"
	"github.com/spf13/cobra"
)

var (
	configFlags *config.FlightDeckConfig

	rootCmd = &cobra.Command{
		Use:   "flightdeck",
		Short: "Flightdeck lets you track flights in the terminal.",
		Long:  "A terminal application for tracking flight and airport data in real time.",
		Run:   run,
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// Runs flightdeck with provided configuration
func run(cmd *cobra.Command, args []string) {
	internal.Init(configFlags)
}

// initialize the command and setup CLI flag options
func init() {
	configFlags = config.NewFlightDeckConfig()

	// set an airport from the CLI
	rootCmd.Flags().StringVarP(
		&configFlags.Airport,
		"airport", "a",
		"",
		"Specify the airport to display",
	)
}
