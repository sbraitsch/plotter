/*
Copyright © 2025 Simon Braitsch <S.Braitsch@gmx.de>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/sbraitsch/plotter/internal/optimizer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "plotter",
	Short: "Attempts to optimize WoW Housing Assignments for maximum satisfaction.",
	Long: `Attempts to optimize WoW Housing Assignments for maximum satisfaction.
	Allows to prioritize neighbors vs. position.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("reading config: %v", err)
		}

		var config optimizer.Config
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("unmarshaling config: %v", err)
		}

		// optimizer.PrettyPrint(config)

		optimizer.Optimize(config)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&cfgFile, "input", "i", "generated.yaml", "path to YAML config file")
}
