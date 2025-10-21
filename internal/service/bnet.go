package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/sbraitsch/plotter/internal/model"
	"github.com/sbraitsch/plotter/internal/storage"
)

type BnetService interface {
	GetProfile(ctx context.Context) (*model.WowProfile, error)
	GetGuildRoster(ctx context.Context, community *model.Community) (*model.Roster, error)
	GetUserGuilds(ctx context.Context) ([]model.Community, error)
}

type bnetServiceImpl struct {
	client  *http.Client
	storage *storage.StorageClient
}

func NewBnetService(client *http.Client, storage *storage.StorageClient) BnetService {
	return &bnetServiceImpl{client: client, storage: storage}
}

func (s *bnetServiceImpl) GetProfile(ctx context.Context) (*model.WowProfile, error) {
	resp, err := s.client.Get("https://eu.api.blizzard.com/profile/user/wow?namespace=profile-eu&locale=en_US")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile model.WowProfile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *bnetServiceImpl) GetGuildRoster(ctx context.Context, community *model.Community) (*model.Roster, error) {
	guildSlug := strings.ToLower(strings.ReplaceAll(community.Name, " ", "-"))
	rosterURL := fmt.Sprintf(
		"https://eu.api.blizzard.com/data/wow/guild/%s/%s/roster?namespace=profile-eu&locale=en_US",
		community.Realm, guildSlug,
	)

	resp, err := s.client.Get(rosterURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch guild roster: %w", err)
	}
	defer resp.Body.Close()

	var roster model.Roster
	if err := json.NewDecoder(resp.Body).Decode(&roster); err != nil {
		return nil, fmt.Errorf("failed to parse guild roster: %w", err)
	}
	return &roster, nil
}

func (s *bnetServiceImpl) GetUserGuilds(ctx context.Context) ([]model.Community, error) {
	profile, err := s.GetProfile(ctx)
	if err != nil {
		log.Printf("Failed to fetch profile as guild data prerequisite: %v", err)
		return nil, err
	}

	guilds, err := getUniqueGuilds(profile, s.client)
	if err != nil {
		log.Printf("Failed to fetch unique guilds: %v", err)
		return nil, err
	}
	saved, err := s.storage.InsertGuilds(ctx, guilds)
	if err != nil {
		log.Printf("Failed to insert guilds into database: %v", err)
		return nil, err
	}
	return saved, nil
}

func getUniqueGuilds(profile *model.WowProfile, client *http.Client) ([]model.Community, error) {
	type result struct {
		model.Community
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
			go func(c model.CharacterResponseSimple) {
				defer wg.Done()

				charUrl := fmt.Sprintf(
					"https://eu.api.blizzard.com/profile/wow/character/%s/%s?namespace=profile-eu&locale=en_US",
					c.Realm.Slug,
					strings.ToLower(c.Name),
				)

				resp, err := client.Get(charUrl)
				if err != nil {
					results <- result{model.Community{}, err}
					return
				}
				defer resp.Body.Close()

				var detail model.CharacterResponseDetailed
				if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
					results <- result{model.Community{}, err}
					return
				}

				if detail.Guild.Name != "" {
					results <- result{model.Community{Id: "", Name: detail.Guild.Name, Realm: c.Realm.Slug}, nil}
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
		guildSet[r.Name] = r.Realm
	}

	// convert set to slice
	guilds := make([]model.Community, 0, len(guildSet))
	for g := range guildSet {
		info := model.Community{
			Name:  g,
			Realm: guildSet[g],
		}
		guilds = append(guilds, info)
	}

	return guilds, nil
}
