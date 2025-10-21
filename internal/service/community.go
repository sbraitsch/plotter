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
	JoinCommunity(ctx context.Context, communityId string) error
	ToggleCommunityLock(ctx context.Context, user *model.User) ([]model.Assignment, error)
	GetAssignments(ctx context.Context, communityId string) ([]model.Assignment, error)
	SetOfficerRank(ctx context.Context, communityId string, officerRank int) error
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

func (s *communityServiceImpl) JoinCommunity(ctx context.Context, communityId string) error {
	user := ctx.Value(middleware.CtxUser).(*model.User)
	community, err := s.storage.GetCommunity(ctx, communityId)
	if err != nil {
		log.Printf("Error retrieving community to join from database: %v", err)
	}

	client := oauth.GetClient(ctx)
	bnetService := NewBnetService(client, s.storage)

	profile, err := bnetService.GetProfile(ctx)
	if err != nil {
		log.Printf("Failed to retrieve wow profile: %v", err)
		return err
	}
	roster, err := bnetService.GetGuildRoster(ctx, community)

	err = s.storage.JoinCommunity(ctx, user, communityId, profile, roster)
	return err
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

func (s *communityServiceImpl) SetOfficerRank(ctx context.Context, communityId string, officerRank int) error {
	return s.storage.SetOfficerRank(ctx, communityId, officerRank)
}
