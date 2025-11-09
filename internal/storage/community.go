package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/sbraitsch/plotter/internal/model"
)

func (s *StorageClient) FinalizeCommunity(ctx context.Context, communityId string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE communities
				 SET finalized = NOT finalized
				 WHERE id = $1`,
		communityId,
	)

	if err != nil {
		log.Printf("Failed to update community values for community %s:%v", communityId, err)
		return fmt.Errorf("Information could not be persisted.")
	}
	return nil
}

func (s *StorageClient) GetCommunityData(ctx context.Context, user *model.User) (*model.CommunityData, error) {
	rows, err := s.db.Query(ctx, `
        SELECT u.battletag, u.char, pm.plot_id, pm.priority
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
		var btag, char string
		var fromNum, toNum *int
		if err := rows.Scan(&btag, &char, &fromNum, &toNum); err != nil {
			return nil, err
		}

		if _, exists := playerMap[btag]; !exists {
			playerMap[btag] = &model.MemberData{
				BattleTag: btag,
				Character: char,
				PlotData:  make(map[int]int),
			}
		}

		if fromNum != nil && toNum != nil {
			playerMap[btag].PlotData[*fromNum] = *toNum
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

func (s *StorageClient) GetFullCommunityData(ctx context.Context, user *model.User) (*model.FullCommunityData, error) {
	rows, err := s.db.Query(ctx, `
		SELECT
			u.battletag,
			u.char,
			COALESCE(u.note, '') AS note,
			a.plot_id,
			a.plot_score,
			pm.plot_id AS mapping_plot_id,
			pm.priority
		FROM users u
		LEFT JOIN assignments a ON a.battletag = u.battletag
		LEFT JOIN plot_mappings pm ON pm.battletag = u.battletag
		WHERE u.community_id = $1
		ORDER BY u.battletag, pm.plot_id
	`, user.Community.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberMap := make(map[string]*model.FullMemberData)

	for rows.Next() {
		var (
			btag, char, note          string
			assignPlotID, assignScore sql.NullInt32
			mappingPlotID, priority   sql.NullInt32
		)

		if err := rows.Scan(&btag, &char, &note, &assignPlotID, &assignScore, &mappingPlotID, &priority); err != nil {
			return nil, err
		}

		member, exists := memberMap[btag]
		if !exists {
			member = &model.FullMemberData{
				Assignment: model.Assignment{
					Battletag: btag,
					Character: char,
				},
				Note:     note,
				PlotData: make(map[int]int),
			}

			// Fill assignment if available
			if assignPlotID.Valid {
				member.Assignment.Plot = int(assignPlotID.Int32)
			}
			if assignScore.Valid {
				member.Assignment.Score = int(assignScore.Int32)
			}

			memberMap[btag] = member
		}

		if mappingPlotID.Valid && priority.Valid {
			member.PlotData[int(mappingPlotID.Int32)] = int(priority.Int32)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	members := make([]model.FullMemberData, 0, len(memberMap))
	for _, m := range memberMap {
		members = append(members, *m)
	}

	return &model.FullCommunityData{
		Id:      user.Community.Id,
		Members: members,
	}, nil
}

func (s *StorageClient) GetCommunity(ctx context.Context, communityId string) (*model.Community, int, error) {
	var community model.Community
	requiredRank := 0
	err := s.db.QueryRow(ctx, `SELECT id, name, realm, member_rank FROM communities WHERE id = $1`, communityId).
		Scan(&community.Id, &community.Name, &community.Realm, &requiredRank)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get community info: %w", err)
	}
	return &community, requiredRank, nil
}

func (s *StorageClient) GetCommunitySize(ctx context.Context, communityId string) (int, error) {
	var count int

	err := s.db.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM users
        WHERE community_id = $1
    `, communityId).Scan(&count)

	if err != nil {
		return 1000, fmt.Errorf("Failed to get community usage: %v", err)
	}
	return count, nil
}

func (s *StorageClient) JoinCommunity(
	ctx context.Context,
	user *model.User,
	requiredRank int,
	communityId string,
	profile *model.WowProfile,
	roster *model.Roster,
) (string, error) {
	minRank := math.MaxInt
	var charName string
	for _, acc := range profile.WowAccounts {
		for _, char := range acc.Characters {
			for _, member := range roster.Members {
				if strings.EqualFold(char.Name, member.Character.Name) {
					if member.Rank < minRank {
						minRank = member.Rank
						charName = member.Character.Name
					}
				}
			}
		}
	}

	if requiredRank < minRank {
		log.Printf("Failed to join community. Rank requirement not fulfilled.")
		return "", fmt.Errorf("Rank requirement not fulfilled.")
	}

	_, err := s.db.Exec(ctx,
		`UPDATE users
			 SET char = $1, community_id = $2, community_rank = $3
			 WHERE battletag = $4`,
		charName, communityId, minRank, user.Battletag,
	)

	if err != nil {
		log.Printf("Failed to update community values for user %v:%v", user, err)
		return "", fmt.Errorf("Information could not be persisted.")
	}

	return charName, nil
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

	_, err = s.db.Exec(ctx,
		`DELETE FROM assignments WHERE community_id=$1`,
		communityId,
	)
	if err != nil {
		return fmt.Errorf("failed to clean up assignments before update: %w", err)
	}

	sqlStr := `INSERT INTO assignments (battletag, char, community_id, plot_id, plot_score) VALUES `
	args := []any{}

	for i, a := range assignments {
		idx := i * 5
		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", idx+1, idx+2, idx+3, idx+4, idx+5)
		args = append(args, a.Battletag, a.Character, communityId, a.Plot, a.Score)
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
		SELECT battletag, char, plot_id, plot_score
		FROM assignments as a
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
		if err := rows.Scan(&a.Battletag, &a.Character, &a.Plot, &a.Score); err != nil {
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

func (s *StorageClient) SetAssignment(ctx context.Context, req *model.SingleAssignmentRequest, communityId string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin assignment transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = s.RegisterManualUsers(ctx, []model.Assignment{model.Assignment{
		Character: req.Char,
		Battletag: req.Battletag,
		Plot:      req.PlotId,
		Score:     0,
	}}, communityId)

	if err != nil {
		log.Printf("Failed to register new community member %s: %v", req.Battletag, err)
		return err
	}

	_, err = s.db.Exec(ctx, `
			DELETE FROM assignments
			WHERE plot_id = $1 AND community_id = $2
		`, req.PlotId, communityId)
	if err != nil {
		log.Printf("Failed to remove plot assignment for %s: %v", req.Battletag, err)
		return err
	}

	_, err = s.db.Exec(ctx, `
		INSERT INTO assignments (battletag, plot_id, char, community_id, plot_score)
		VALUES ($1, $2, $3, $4, 0)
		ON CONFLICT (battletag)
		DO UPDATE SET
			plot_id = EXCLUDED.plot_id,
			char = EXCLUDED.char,
			plot_score = 0
`, req.Battletag, req.PlotId, req.Char, communityId)
	if err != nil {
		log.Printf("Failed to update plot assignment for %s to %d: %v", req.Battletag, req.PlotId, err)
		return err
	}

	return nil
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

func (s *StorageClient) SetOfficerRank(ctx context.Context, communityId string, req *model.CommunityRankRequest) error {
	_, err := s.db.Exec(ctx, `
			UPDATE communities
			SET officer_rank = $1, member_rank = $2
			WHERE id = $3::uuid
		`, req.AdminRank, req.MemberRank, communityId)
	if err != nil {
		log.Printf("Failed to update community %s's rank settings: %v", communityId, err)
		return err
	}

	return nil
}

func (s *StorageClient) GetCommunitySettings(ctx context.Context, communityId string) (*model.Settings, error) {
	var officerRank, memberRank sql.NullInt32
	err := s.db.QueryRow(ctx,
		`SELECT officer_rank, member_rank
			     FROM communities
				 WHERE id = $1`,
		communityId,
	).Scan(&officerRank, &memberRank)

	if err != nil {
		log.Printf("Failed to retrieve settings for community %s: %v", communityId, err)
		return nil, err
	}

	return &model.Settings{OfficerRank: int(officerRank.Int32), MemberRank: int(memberRank.Int32)}, nil
}
