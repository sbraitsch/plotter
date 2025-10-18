package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sbraitsch/plotter/internal/middleware"
	"github.com/sbraitsch/plotter/internal/service/oauth"
)

type Community struct {
	Id      string       `json:"id"`
	Members []PlayerData `json:"members"`
}

type PlayerData struct {
	BattleTag string      `json:"battletag"`
	PlotData  map[int]int `json:"plotData"`
}

func GetCommunity(ctx context.Context, db *pgxpool.Pool) (*Community, error) {
	user, ok := ctx.Value(middleware.CtxUser).(middleware.UserContext)
	if !ok || len(user.CommunityId) == 0 {
		return nil, fmt.Errorf("community not found in context")
	}

	rows, err := db.Query(ctx, `
        SELECT u.battletag, pm.plot_id, pm.priority
        FROM users u
        LEFT JOIN plot_mappings pm ON pm.battletag = u.battletag
		WHERE u.community_id = $1
        ORDER BY u.battletag, pm.plot_id
    `, user.CommunityId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	playerMap := make(map[string]*PlayerData)

	for rows.Next() {
		var name string
		var fromNum, toNum *int
		if err := rows.Scan(&name, &fromNum, &toNum); err != nil {
			return nil, err
		}

		if _, exists := playerMap[name]; !exists {
			playerMap[name] = &PlayerData{
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
	members := make([]PlayerData, 0, len(playerMap))
	for _, pd := range playerMap {
		members = append(members, *pd)
	}

	community := &Community{Id: user.CommunityId, Members: members}
	return community, nil
}

func JoinCommunity(ctx context.Context, db *pgxpool.Pool, communityId string) error {
	user := ctx.Value(middleware.CtxUser).(middleware.UserContext)
	var guildName, realmSlug string
	err := db.QueryRow(ctx, `SELECT name, realm FROM communities WHERE id = $1`, communityId).Scan(&guildName, &realmSlug)
	if err != nil {
		return fmt.Errorf("failed to get community info: %w", err)
	}

	client := oauth.GetClient(ctx)
	profile, err := GetProfile(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	guildSlug := strings.ToLower(strings.ReplaceAll(guildName, " ", "-"))
	rosterURL := fmt.Sprintf(
		"https://eu.api.blizzard.com/data/wow/guild/%s/%s/roster?namespace=profile-eu&locale=en_US",
		realmSlug, guildSlug,
	)

	resp, err := client.Get(rosterURL)
	if err != nil {
		return fmt.Errorf("failed to fetch guild roster: %w", err)
	}
	defer resp.Body.Close()

	var roster struct {
		Members []struct {
			Character struct {
				Name  string `json:"name"`
				Realm struct {
					Slug string `json:"slug"`
				} `json:"realm"`
			} `json:"character"`
			Rank int `json:"rank"`
		} `json:"members"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&roster); err != nil {
		return fmt.Errorf("failed to parse guild roster: %w", err)
	}

	minRank := math.MaxInt
	for _, acc := range profile.WowAccounts {
		for _, char := range acc.Characters {
			for _, member := range roster.Members {
				if strings.EqualFold(char.Name, member.Character.Name) &&
					strings.EqualFold(char.Realm.Slug, member.Character.Realm.Slug) {
					if member.Rank < minRank {
						minRank = member.Rank
					}
				}
			}
		}
	}

	_, err = db.Exec(ctx,
		`UPDATE users
		 SET community_id = $1, community_rank = $2
		 WHERE battletag = $3`,
		communityId, minRank, user.Battletag,
	)

	if err != nil {
		return err
	}

	return nil
}
