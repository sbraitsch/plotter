package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const adminName = "admin"

func SeedAdminPlayers(ctx context.Context, pool *pgxpool.Pool) ([]string, error) {
	// 1️⃣ Fetch all existing admins
	rows, err := pool.Query(ctx, `SELECT uuid FROM players WHERE is_admin = TRUE`)
	if err != nil {
		return nil, fmt.Errorf("failed to query admin players: %w", err)
	}
	defer rows.Close()

	var adminUUIDs []string
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			return nil, fmt.Errorf("failed to scan admin uuid: %w", err)
		}
		adminUUIDs = append(adminUUIDs, uuid)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating admin rows: %w", err)
	}

	// 2️⃣ If we already have admins, return them
	if len(adminUUIDs) > 0 {
		log.Printf("✅ Found %d admin user(s): %v", len(adminUUIDs), adminUUIDs)
		return adminUUIDs, nil
	}

	// 3️⃣ Otherwise, create the initial admin
	var uuid string
	err = pool.QueryRow(ctx,
		`INSERT INTO players (name, is_admin) VALUES ($1, TRUE) RETURNING uuid`,
		adminName,
	).Scan(&uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial admin: %w", err)
	}

	log.Printf("✅ Created initial admin user '%s' with UUID: %s", adminName, uuid)

	return []string{uuid}, nil
}
