package service

import (
	"context"
	"fmt"
	"log"

	"github.com/sbraitsch/plotter/internal/middleware"
	"github.com/sbraitsch/plotter/internal/model"
	"github.com/sbraitsch/plotter/internal/service/oauth"
	"github.com/sbraitsch/plotter/internal/storage"
)

type CommunityService interface {
	GetCommunityData(ctx context.Context) (*model.CommunityData, error)
	JoinCommunity(ctx context.Context, communityId string) (string, error)
	ToggleCommunityLock(ctx context.Context, user *model.User) ([]model.Assignment, error)
	GetAssignments(ctx context.Context, communityId string) ([]model.Assignment, error)
	SetCommunitySettings(ctx context.Context, communityId string, req *model.CommunityRankRequest) error
	GetCommunitySettings(ctx context.Context, communityId string) (*model.Settings, error)
	DownloadCommunityData(ctx context.Context) (*model.FullCommunityData, error)
	UploadCommunityData(ctx context.Context, data *model.AssignmentUpload) ([]model.Assignment, error)
}

type communityServiceImpl struct {
	storage *storage.StorageClient
}

func NewCommunityService(storage *storage.StorageClient) CommunityService {
	return &communityServiceImpl{storage: storage}
}

func (s *communityServiceImpl) GetCommunityData(ctx context.Context) (*model.CommunityData, error) {
	user, ok := ctx.Value(middleware.CtxUser).(*model.User)
	if !ok || len(user.Community.Id) == 0 {
		return nil, fmt.Errorf("community not found in context")
	}
	community, err := s.storage.GetCommunityData(ctx, user)
	if err != nil {
		log.Printf("Failed to retrieve community data from database: %v", err)
		return nil, err
	}
	return community, nil
}

func (s *communityServiceImpl) DownloadCommunityData(ctx context.Context) (*model.FullCommunityData, error) {
	user, ok := ctx.Value(middleware.CtxUser).(*model.User)
	if !ok || len(user.Community.Id) == 0 {
		return nil, fmt.Errorf("community not found in context")
	}
	community, err := s.storage.GetFullCommunityData(ctx, user)
	if err != nil {
		log.Printf("Failed to retrieve community data from database: %v", err)
		return nil, err
	}
	return community, nil
}

func (s *communityServiceImpl) UploadCommunityData(ctx context.Context, data *model.AssignmentUpload) ([]model.Assignment, error) {
	user := ctx.Value(middleware.CtxUser).(*model.User)
	assignments := make([]model.Assignment, 0, len(data.Members))
	for _, member := range data.Members {
		if member.Assignment.Battletag == "" {
			continue
		}
		assignments = append(assignments, model.Assignment{
			Battletag: member.Assignment.Battletag,
			Score:     member.Assignment.Score,
			Character: member.Assignment.Character,
			Plot:      member.Assignment.Plot,
		})
	}

	if len(assignments) > 53 {
		return nil, fmt.Errorf("too many assignments. community is overcrowded.")
	}
	err := s.storage.RegisterManualUsers(ctx, assignments, user.Community.Id)
	err = s.storage.PersistAndLock(ctx, assignments, user.Community.Id)
	if err != nil {
		log.Printf("Error persisting overwritten assignments: %v", err)
		return nil, err
	}
	log.Printf("Community %s locked.", user.Community.Id)
	return assignments, nil
}

func (s *communityServiceImpl) JoinCommunity(ctx context.Context, communityId string) (string, error) {
	user := ctx.Value(middleware.CtxUser).(*model.User)
	occupancy, err := s.storage.GetCommunitySize(ctx, communityId)
	if err != nil {
		log.Printf("Error retrieving community occupancy from database: %v", err)
		return "", err
	}
	if occupancy > 53 {
		return "", fmt.Errorf("Community is full. Apologies.")
	}

	community, requiredRank, err := s.storage.GetCommunity(ctx, communityId)
	if err != nil {
		log.Printf("Error retrieving community to join from database: %v", err)
		return "", err
	}

	client := oauth.GetClient(ctx)
	bnetService := NewBnetService(client, s.storage)

	profile, err := bnetService.GetProfile(ctx)
	if err != nil {
		log.Printf("Failed to retrieve wow profile: %v", err)
		return "", err
	}
	roster, err := bnetService.GetGuildRoster(ctx, community)

	joinedChar, err := s.storage.JoinCommunity(ctx, user, requiredRank, communityId, profile, roster)
	return joinedChar, err
}

func (s *communityServiceImpl) ToggleCommunityLock(ctx context.Context, user *model.User) ([]model.Assignment, error) {
	// unlock if locked
	if user.Community.Locked {
		err := s.storage.UnlockCommunity(ctx, user.Community.Id)
		if err != nil {
			log.Printf("Failed to unlock community: %v", err)
			return nil, err
		}
		log.Printf("Community %s locked.", user.Community.Id)
		return nil, nil
	}
	community, err := s.GetCommunityData(ctx)
	if err != nil {
		log.Printf("Failed to fetch community to optimize: %v", err)
		return nil, err
	}

	assignments := community.Optimize()

	err = s.storage.PersistAndLock(ctx, assignments, community.Id)
	if err != nil {
		log.Printf("Error persisting assignments: %v", err)
		return nil, err
	}
	log.Printf("Community %s locked.", community.Id)
	return assignments, nil
}

func (s *communityServiceImpl) GetAssignments(ctx context.Context, communityId string) ([]model.Assignment, error) {
	return s.storage.GetAssignments(ctx, communityId)
}

func (s *communityServiceImpl) SetCommunitySettings(ctx context.Context, communityId string, req *model.CommunityRankRequest) error {
	return s.storage.SetOfficerRank(ctx, communityId, req)
}

func (s *communityServiceImpl) GetCommunitySettings(ctx context.Context, communityId string) (*model.Settings, error) {
	return s.storage.GetCommunitySettings(ctx, communityId)
}
