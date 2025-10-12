package cmd

import (
	"fmt"

	"github.com/sbraitsch/plotter/internal/optimizer"
	"github.com/spf13/cobra"
)

var url string

var sheetCmd = &cobra.Command{
	Use:   "sheet",
	Short: "Read data from Google Sheet",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := optimizer.ReadConfigCSV(url)
		if err != nil {
			fmt.Printf("Error parsing CSV: %v", err)
		}

		optimizer.Optimize(cfg)
	},
}

func init() {
	rootCmd.AddCommand(sheetCmd)

	// Flags
	sheetCmd.Flags().StringVarP(
		&url,
		"sheet",
		"s",
		"https://docs.google.com/spreadsheets/d/e/2PACX-1vR_qx5QOAwTZZKpl-FXr4PumRsieftY36a4Dv8Vn2uBvDWlG8Z7Q50OKjWiyu6oao7zOZ5ASjEkh5Vz/pub?gid=0&single=true&output=csv",
		"URL to Google Sheet",
	)
}
