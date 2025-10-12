/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sbraitsch/plotter/internal/api"
	"github.com/sbraitsch/plotter/internal/config"
	"github.com/sbraitsch/plotter/internal/db"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the plotter functions to the web",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()
		cfg := config.Load()
		ctx := context.Background()

		pool := db.Connect(ctx, cfg.DBURL)

		srv := api.Server{DB: pool}

		addr := fmt.Sprintf(":%s", cfg.Port)

		log.Printf("Starting server on %s\n", addr)
		if err := srv.Start(addr); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
