/*
Copyright © 2025 Simon Braitsch <S.Braitsch@gmx.de>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "plotter",
	Short: "Attempts to optimize WoW Housing Assignments for maximum satisfaction.",
	Long: `Attempts to optimize WoW Housing Assignments for maximum satisfaction.
	Allows to prioritize neighbors vs. position.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
