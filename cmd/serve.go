/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
		defer pool.Close()

		db.RunMigrations(cfg.DBURL)
		adminUUID, err := db.SeedAdminPlayer(ctx, pool)
		if err != nil {
			log.Fatalf("failed to seed admin: %v", err)
		}
		log.Printf("Admin UUID: %s", adminUUID)

		srv := api.Server{DB: pool, AdminUUID: adminUUID}
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
