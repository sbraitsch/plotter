package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const adminName = "Tejb"

func SeedAdminPlayer(ctx context.Context, pool *pgxpool.Pool) (string, error) {
	var uuid string
	err := pool.QueryRow(ctx, `SELECT uuid FROM players WHERE name=$1`, adminName).Scan(&uuid)
	if err == nil {
		log.Printf("Admin user '%s' already exists with UUID: %s", adminName, uuid)
		return uuid, nil
	}

	// create admin
	err = pool.QueryRow(ctx, `INSERT INTO players (name) VALUES ($1) RETURNING uuid`, adminName).Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("failed to create admin: %w", err)
	}

	log.Printf("✅ Created admin user '%s' with UUID: %s", adminName, uuid)
	return uuid, nil
}
