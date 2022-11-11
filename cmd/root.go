package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stopwatch2",
	Short: "StopWatch Bike & Run races",
}

var debug bool
var force bool

const TimeFormat = "2006-01-02:15:04:05"

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(newImportTeamsCmd())
	rootCmd.AddCommand(newShellCmd())
	rootCmd.AddCommand(newGenerateReportCmd())
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Display debug logs")
	rootCmd.PersistentFlags().BoolVar(&force, "force", false, "Force-write in output file even if it exists (existing content will be lost)")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
