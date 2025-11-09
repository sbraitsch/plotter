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
		battletag, char, note, communityName, communityID, realm, accessToken sql.NullString
		officerRank, communityRank                                            sql.NullInt32
		locked, finalized                                                     sql.NullBool
		expiry                                                                sql.NullTime
	)

	err := s.db.QueryRow(ctx,
		`SELECT
			u.battletag,
			u.char,
			u.note,
			u.community_id,
			c.name AS community_name,
			c.officer_rank,
			c.locked,
			c.finalized,
			c.realm,
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
		&note,
		&communityID,
		&communityName,
		&officerRank,
		&locked,
		&finalized,
		&realm,
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
		Note:      note.String,
		Community: model.UserCommunity{
			Id:          communityID.String,
			Name:        communityName.String,
			OfficerRank: int(officerRank.Int32),
			Locked:      locked.Bool,
			Realm:       realm.String,
			Finalized:   finalized.Bool,
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

func (s *StorageClient) RegisterManualUsers(ctx context.Context, assignments []model.Assignment, communityID string) error {
	sessionToken := uuid.New().String()

	var memberRank int
	err := s.db.QueryRow(ctx, `SELECT member_rank FROM communities WHERE id = $1`, communityID).Scan(&memberRank)
	if err != nil {
		return fmt.Errorf("failed to fetch community member_rank: %w", err)
	}
	for _, a := range assignments {
		if a.Battletag == "" {
			continue
		}
		_, err := s.db.Exec(ctx, `
			INSERT INTO users (battletag, char, community_id, community_rank, access_token, expiry, session_id)
			VALUES ($1, $2, $3, $4, gen_random_uuid()::text, NOW() + INTERVAL '24 hours', $5)
			ON CONFLICT (battletag) DO NOTHING
		`, a.Battletag, a.Character, communityID, memberRank, sessionToken)

		if err != nil {
			log.Printf("Failed to insert user %s: %v", a.Battletag, err)
			return err
		}
	}
	return nil
}

func (s *StorageClient) SetNote(ctx context.Context, user *model.User, note string) error {
	_, err := s.db.Exec(ctx, `UPDATE users SET note=$1 WHERE battletag=$2`, note, user.Battletag)
	if err != nil {
		log.Printf("failed to set user note: %v", err)
		return err
	}
	return nil
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
