package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sbraitsch/plotter/internal/middleware"
)

func UpdatePlayerData(ctx context.Context, db *pgxpool.Pool, mapping map[int]int) error {
	user := ctx.Value(middleware.CtxUser).(middleware.UserContext)

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	plotIDs := make([]any, 0, len(mapping))
	for plotID := range mapping {
		plotIDs = append(plotIDs, plotID)
	}

	if len(plotIDs) > 0 {
		placeholders := make([]string, 0, len(mapping))
		args := []any{user.Battletag}
		idx := 2
		for plotID := range mapping {
			placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
			args = append(args, plotID)
			idx++
		}

		query := fmt.Sprintf(
			`DELETE FROM plot_mappings WHERE battletag=$1 AND plot_id NOT IN (%s)`,
			strings.Join(placeholders, ","),
		)
		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			log.Fatalf("failed to remove mappings: %v", err)
			return err
		}
	} else {
		_, err := tx.Exec(ctx, `DELETE FROM plot_mappings WHERE battletag=$1`, user.Battletag)
		if err != nil {
			return err
		}
	}

	for plotId, priority := range mapping {
		_, err = tx.Exec(ctx, `
			INSERT INTO plot_mappings (battletag, plot_id, priority)
			VALUES ($1, $2, $3)
			ON CONFLICT (battletag, plot_id)
			DO UPDATE SET priority = EXCLUDED.priority
		`, user.Battletag, plotId, priority)
		if err != nil {
			log.Fatalf("failed to save mapping %d → %d: %v", plotId, priority, err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
