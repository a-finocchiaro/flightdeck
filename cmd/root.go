package cmd

import (
	"github.com/a-finocchiaro/flightdeck/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flightdeck",
	Short: "Flightdeck lets you track flights in the terminal.",
	Long:  "A terminal application for tracking flight and airport data in real time.",
	Run:   run,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func run(cmd *cobra.Command, args []string) {
	internal.Init()
}
