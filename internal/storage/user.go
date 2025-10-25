package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/sbraitsch/plotter/internal/model"
	"golang.org/x/oauth2"
)

func (s *StorageClient) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	var (
		battletag, char, communityName, communityID, accessToken sql.NullString
		officerRank, communityRank                               sql.NullInt32
		locked                                                   sql.NullBool
		expiry                                                   sql.NullTime
	)

	err := s.db.QueryRow(ctx,
		`SELECT
			u.battletag,
			u.char,
			u.community_id,
			c.name AS community_name,
			c.officer_rank,
			c.locked,
			u.community_rank,
			u.access_token,
			u.expiry
		FROM users u
		LEFT JOIN communities c
			ON u.community_id = c.id
		WHERE u.session_id = $1`,
		token,
	).Scan(
		&battletag,
		&char,
		&communityID,
		&communityName,
		&officerRank,
		&locked,
		&communityRank,
		&accessToken,
		&expiry,
	)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Battletag: battletag.String,
		Char:      char.String,
		Community: model.UserCommunity{
			Id:          communityID.String,
			Name:        communityName.String,
			OfficerRank: int(officerRank.Int32),
			Locked:      locked.Bool,
		},
		CommunityRank: int(communityRank.Int32),
		AccessToken:   accessToken.String,
		Expiry:        expiry.Time,
	}

	return user, nil
}

func (s *StorageClient) RegisterUser(ctx context.Context, battletag string, token *oauth2.Token) (string, error) {
	sessionToken := uuid.New().String()
	_, err := s.db.Exec(ctx, `INSERT INTO users(battletag, session_id, access_token, expiry)
                      VALUES($1, $2, $3, $4)
                      ON CONFLICT(battletag) DO UPDATE
                      SET session_id=$2, access_token=$3, expiry=$4`,
		battletag, sessionToken, token.AccessToken, token.Expiry)

	if err != nil {
		log.Printf("Failed to insert new user: %v", err)
		return "", err
	}
	return sessionToken, nil
}

func (s *StorageClient) SavePlotMappings(ctx context.Context, user *model.User, mappings map[int]int) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	plotIDs := make([]any, 0, len(mappings))
	for plotID := range mappings {
		plotIDs = append(plotIDs, plotID)
	}

	if len(plotIDs) > 0 {
		placeholders := make([]string, 0, len(mappings))
		args := []any{user.Battletag}
		idx := 2
		for plotID := range mappings {
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
			log.Printf("failed to remove mappings: %v", err)
			return err
		}
	} else {
		_, err := tx.Exec(ctx, `DELETE FROM plot_mappings WHERE battletag=$1`, user.Battletag)
		if err != nil {
			return err
		}
	}

	for plotId, priority := range mappings {
		_, err = tx.Exec(ctx, `
			INSERT INTO plot_mappings (battletag, plot_id, priority)
			VALUES ($1, $2, $3)
			ON CONFLICT (battletag, plot_id)
			DO UPDATE SET priority = EXCLUDED.priority
		`, user.Battletag, plotId, priority)
		if err != nil {
			log.Printf("failed to save mapping %d → %d: %v", plotId, priority, err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
