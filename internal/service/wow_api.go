package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sbraitsch/plotter/internal/models"
	"github.com/sbraitsch/plotter/internal/service/oauth"
)

type GuildInfo struct {
	Name      string
	RealmSlug string
}

func GetGuilds(ctx context.Context, db *pgxpool.Pool) ([]CommunityInfo, error) {
	client := oauth.GetClient(ctx)
	profile, err := GetProfile(ctx, client)
	if err != nil {
		return nil, err
	}

	guilds, err := getUniqueGuilds(profile, client)
	if err != nil {
		return nil, err
	}
	saved, err := insertGuilds(ctx, db, guilds)
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func GetProfile(ctx context.Context, client *http.Client) (*models.ProfileResponse, error) {
	resp, err := client.Get("https://eu.api.blizzard.com/profile/user/wow?namespace=profile-eu&locale=en_US")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile models.ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func getUniqueGuilds(profile *models.ProfileResponse, client *http.Client) ([]GuildInfo, error) {
	type result struct {
		GuildInfo
		Err error
	}

	results := make(chan result, 100)
	var wg sync.WaitGroup

	for _, acc := range profile.WowAccounts {
		for _, char := range acc.Characters {
			if char.Level < 80 {
				continue
			}
			wg.Add(1)
			go func(c models.CharacterResponseSimple) {
				defer wg.Done()

				charUrl := fmt.Sprintf(
					"https://eu.api.blizzard.com/profile/wow/character/%s/%s?namespace=profile-eu&locale=en_US",
					c.Realm.Slug,
					strings.ToLower(c.Name),
				)

				resp, err := client.Get(charUrl)
				if err != nil {
					results <- result{GuildInfo{}, err}
					return
				}
				defer resp.Body.Close()

				var detail models.CharacterResponseDetailed
				if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
					results <- result{GuildInfo{}, err}
					return
				}

				if detail.Guild.Name != "" {
					results <- result{GuildInfo{detail.Guild.Name, c.Realm.Slug}, nil}
				}
			}(char)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	guildSet := make(map[string]string)
	for r := range results {
		if r.Err != nil {
			return nil, r.Err
		}
		guildSet[r.Name] = r.RealmSlug
	}

	// convert set to slice
	guilds := make([]GuildInfo, 0, len(guildSet))
	for g := range guildSet {
		info := GuildInfo{
			Name:      g,
			RealmSlug: guildSet[g],
		}
		guilds = append(guilds, info)
	}

	return guilds, nil
}

func insertGuilds(ctx context.Context, db *pgxpool.Pool, guilds []GuildInfo) ([]CommunityInfo, error) {
	sqlStr := `INSERT INTO communities (name, realm) VALUES `
	args := []any{}

	for i, g := range guilds {
		idx := i * 2
		sqlStr += fmt.Sprintf("($%d, $%d),", idx+1, idx+2)
		args = append(args, g.Name, g.RealmSlug)
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")
	sqlStr += " ON CONFLICT (name) DO NOTHING"

	_, err := db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(guilds))
	for i, g := range guilds {
		names[i] = g.Name
	}

	rows, err := db.Query(ctx,
		`SELECT id, name
		     FROM communities
			 WHERE name = ANY($1)`,
		names,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var saved []CommunityInfo
	for rows.Next() {
		var c CommunityInfo
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return nil, err
		}
		saved = append(saved, c)
	}

	return saved, nil
}
