package storage

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/sbraitsch/plotter/internal/model"
)

func (s *StorageClient) GetCommunityData(ctx context.Context, user *model.User) (*model.CommunityData, error) {
	rows, err := s.db.Query(ctx, `
        SELECT u.battletag, pm.plot_id, pm.priority
        FROM users u
        LEFT JOIN plot_mappings pm ON pm.battletag = u.battletag
			WHERE u.community_id = $1
        ORDER BY u.battletag, pm.plot_id
    `, user.Community.Id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	playerMap := make(map[string]*model.MemberData)

	for rows.Next() {
		var name string
		var fromNum, toNum *int
		if err := rows.Scan(&name, &fromNum, &toNum); err != nil {
			return nil, err
		}

		if _, exists := playerMap[name]; !exists {
			playerMap[name] = &model.MemberData{
				BattleTag: name,
				PlotData:  make(map[int]int),
			}
		}

		if fromNum != nil && toNum != nil {
			playerMap[name].PlotData[*fromNum] = *toNum
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	members := make([]model.MemberData, 0, len(playerMap))
	for _, pd := range playerMap {
		members = append(members, *pd)
	}

	community := &model.CommunityData{Id: user.Community.Id, Members: members}
	return community, nil
}

func (s *StorageClient) GetCommunity(ctx context.Context, communityId string) (*model.Community, error) {
	var community model.Community
	err := s.db.QueryRow(ctx, `SELECT id, name, realm FROM communities WHERE id = $1`, communityId).Scan(&community.Id, &community.Name, &community.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to get community info: %w", err)
	}
	return &community, nil
}

func (s *StorageClient) JoinCommunity(ctx context.Context, user *model.User, communityId string, profile *model.WowProfile, roster *model.Roster) error {
	minRank := math.MaxInt
	for _, acc := range profile.WowAccounts {
		for _, char := range acc.Characters {
			for _, member := range roster.Members {
				if strings.EqualFold(char.Name, member.Character.Name) {
					if member.Rank < minRank {
						minRank = member.Rank
					}
				}
			}
		}
	}

	_, err := s.db.Exec(ctx,
		`UPDATE users
			 SET community_id = $1, community_rank = $2
			 WHERE battletag = $3`,
		communityId, minRank, user.Battletag,
	)

	if err != nil {
		log.Printf("Failed to update community values for user %v:%v", user, err)
		return err
	}

	return nil
}

func (s *StorageClient) UnlockCommunity(ctx context.Context, communityId string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE communities
		SET locked = false
		WHERE id = $1`,
		communityId,
	)
	return err
}

func (s *StorageClient) PersistAndLock(ctx context.Context, assignments []model.Assignment, communityId string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin lock transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	sqlStr := `INSERT INTO assignments (battletag, community_id, plot_id, plot_score) VALUES `
	args := []any{}

	for i, a := range assignments {
		idx := i * 2
		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d),", idx+1, idx+2, idx+3, idx+4)
		args = append(args, a.Battletag, communityId, a.Plot, a.Score)
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")
	sqlStr += ` ON CONFLICT (battletag)
		        DO UPDATE SET
                  plot_id = EXCLUDED.plot_id,
                  plot_score = EXCLUDED.plot_score`

	_, err = tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE communities
		SET locked = true
		WHERE id = $1`,
		communityId,
	)
	if err = tx.Commit(ctx); err != nil {
		log.Printf("failed to commit lock transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *StorageClient) GetAssignments(ctx context.Context, communityId string) ([]model.Assignment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT battletag, plot_id, plot_score
		FROM assignments
		WHERE community_id = $1
	`, communityId)
	if err != nil {
		log.Printf("Assignment query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	assignments := []model.Assignment{}

	for rows.Next() {
		var a model.Assignment
		if err := rows.Scan(&a.Battletag, &a.Plot, &a.Score); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error reading assignments from database: %v", err)
		return nil, err
	}

	return assignments, nil
}

func (s *StorageClient) InsertGuilds(ctx context.Context, guilds []model.Community) ([]model.Community, error) {
	tx, err := s.db.Begin(ctx)
	sqlStr := `INSERT INTO communities (name, realm) VALUES `
	args := []any{}

	for i, g := range guilds {
		idx := i * 2
		sqlStr += fmt.Sprintf("($%d, $%d),", idx+1, idx+2)
		args = append(args, g.Name, g.Realm)
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")
	sqlStr += " ON CONFLICT (name) DO NOTHING"

	_, err = tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(guilds))
	for i, g := range guilds {
		names[i] = g.Name
	}

	rows, err := tx.Query(ctx,
		`SELECT id, name, realm, locked
		     FROM communities
			 WHERE name = ANY($1)`,
		names,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var saved []model.Community
	for rows.Next() {
		var c model.Community
		if err := rows.Scan(&c.Id, &c.Name, &c.Realm, &c.Locked); err != nil {
			return nil, err
		}
		saved = append(saved, c)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("failed to commit guild insert transaction: %v", err)
		return nil, err
	}

	return saved, nil
}

func (s *StorageClient) SetOfficerRank(ctx context.Context, communityId string, officerRank int) error {
	_, err := s.db.Exec(ctx,
		`UPDATE communities
			 SET officer_rank = $1
			 WHERE id = $2`,
		officerRank, communityId,
	)

	return err
}
