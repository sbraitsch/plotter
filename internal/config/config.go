package config

import (
	"log"
	"os"
)

type Config struct {
	DBURL string
	Port  string
}

func Load() Config {
	cfg := Config{
		DBURL: os.Getenv("DATABASE_URL"),
		Port:  os.Getenv("PORT"),
	}

	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL not set")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}
