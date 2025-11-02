package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/sbraitsch/plotter/internal/middleware"
	"github.com/sbraitsch/plotter/internal/model"
	"github.com/sbraitsch/plotter/internal/service/oauth"
	"github.com/sbraitsch/plotter/internal/storage"
	"golang.org/x/oauth2"
)

type UserService interface {
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
	Validate(ctx context.Context) (*model.ValidatedUser, error)
	RegisterUser(code string, oauth *oauth2.Config) (string, error)
	UpdateMappings(ctx context.Context, mappings map[int]int) (*model.CommunityData, error)
	SetNote(ctx context.Context, note string) error
	ListAvailableCommunities(ctx context.Context) ([]model.Community, error)
}

type userServiceImpl struct {
	storage *storage.StorageClient
}

func NewUserService(storage *storage.StorageClient) UserService {
	return &userServiceImpl{storage: storage}
}

func (s *userServiceImpl) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	user, err := s.storage.GetUserByToken(ctx, token)
	if err != nil {
		log.Printf("Failed to retrieve user from database: %v", err)
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) ListAvailableCommunities(ctx context.Context) ([]model.Community, error) {
	client := oauth.GetClient(ctx)
	bnetService := NewBnetService(client, s.storage)
	return bnetService.GetUserGuilds(ctx)
}

func (s *userServiceImpl) Validate(ctx context.Context) (*model.ValidatedUser, error) {
	user := ctx.Value(middleware.CtxUser).(*model.User)
	return &model.ValidatedUser{
		Battletag: user.Battletag,
		Char:      user.Char,
		Note:      user.Note,
		IsAdmin:   user.CommunityRank <= user.Community.OfficerRank,
		Community: model.ValidatedCommunity{
			Id:     user.Community.Id,
			Name:   user.Community.Name,
			Realm:  user.Community.Realm,
			Locked: user.Community.Locked,
		},
	}, nil
}

func (s *userServiceImpl) RegisterUser(code string, oauth *oauth2.Config) (string, error) {
	ctx := context.Background()

	token, err := oauth.Exchange(ctx, code)
	if err != nil {
		log.Printf("Failed token exchange: %v", err)
		return "", err
	}

	client := oauth.Client(ctx, token)
	resp, err := client.Get("https://oauth.battle.net/oauth/userinfo")
	if err != nil {
		log.Printf("Failed to fetch user profile: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var profile struct {
		Battletag string `json:"battletag"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		log.Printf("Failed to parse user profile: %v", err)
		return "", err
	}

	sessionToken, err := s.storage.RegisterUser(ctx, profile.Battletag, token)
	if err != nil {
		log.Printf("Failed to register new user %v", err)
		return "", err
	}
	return sessionToken, nil
}

func (s *userServiceImpl) UpdateMappings(ctx context.Context, mappings map[int]int) (*model.CommunityData, error) {

	user, ok := ctx.Value(middleware.CtxUser).(*model.User)
	if !ok || len(user.Community.Id) == 0 {
		return nil, fmt.Errorf("community not found in context")
	}

	err := s.storage.SavePlotMappings(ctx, user, mappings)
	if err != nil {
		return nil, err
	}

	community, err := s.storage.GetCommunityData(ctx, user)
	if err != nil {
		log.Printf("Failed to retrieve community data from database: %v", err)
		return nil, err
	}
	return community, nil
}

func (s *userServiceImpl) SetNote(ctx context.Context, note string) error {

	user, ok := ctx.Value(middleware.CtxUser).(*model.User)
	if !ok || len(user.Community.Id) == 0 {
		return fmt.Errorf("community not found in context")
	}

	err := s.storage.SetNote(ctx, user, note)
	if err != nil {
		return err
	}

	return nil
}
