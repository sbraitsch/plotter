package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdatePlayerData(ctx context.Context, db *pgxpool.Pool, name string, mapping map[int]int) error {

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var playerID int
	err = tx.QueryRow(ctx, `SELECT id FROM players WHERE name = $1`, name).Scan(&playerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("player '%s' not found", name)
		}
		return fmt.Errorf("failed to fetch player: %w", err)
	}

	for fromNum, toNum := range mapping {
		_, err = tx.Exec(ctx, `
			INSERT INTO player_mappings (player_id, from_num, to_num)
			VALUES ($1, $2, $3)
			ON CONFLICT (player_id, from_num)
			DO UPDATE SET to_num = EXCLUDED.to_num
		`, playerID, fromNum, toNum)
		if err != nil {
			return fmt.Errorf("failed to save mapping %d → %d: %w", fromNum, toNum, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
