package cmd

import (
	"github.com/sbraitsch/plotter/internal/optimizer"
	"github.com/spf13/cobra"
)

var (
	playerCount int
	plotCount   int
	weightCount int
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a random plot/player config YAML",
	Run: func(cmd *cobra.Command, args []string) {
		optimizer.Generate(playerCount, plotCount, weightCount)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Flags
	generateCmd.Flags().IntVarP(&playerCount, "players", "p", 35, "Number of players")
	generateCmd.Flags().IntVarP(&plotCount, "plots", "l", 50, "Number of plots")
	generateCmd.Flags().IntVarP(&weightCount, "weights", "w", 5, "Number of prioritized plots per player")
}
