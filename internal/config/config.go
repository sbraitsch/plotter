package config

import (
	"log"
	"os"

	"github.com/sbraitsch/plotter/internal/api"
)

func Load() api.Config {
	cfg := api.Config{
		DbUrl:        os.Getenv("DATABASE_URL"),
		Port:         os.Getenv("PORT"),
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	if cfg.DbUrl == "" {
		log.Fatal("DATABASE_URL not set")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}
