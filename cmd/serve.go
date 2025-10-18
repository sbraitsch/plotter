/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

		pool := db.ConnectWithRetry(ctx, cfg.DbUrl, 10, 2*time.Second)
		defer pool.Close()

		db.RunMigrations(cfg.DbUrl)

		srv := api.NewServer(pool, cfg)
		addr := fmt.Sprintf(":%s", cfg.Port)

		log.Printf("Server listening on port %s\n", addr)
		if err := http.ListenAndServe(addr, srv.Router()); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
