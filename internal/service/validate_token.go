package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sbraitsch/plotter/internal/middleware"
)

var ErrInvalidToken = errors.New("invalid token")

type PlayerValidationResponse struct {
	Battletag string        `json:"battletag"`
	IsAdmin   bool          `json:"isAdmin"`
	Community CommunityInfo `json:"community"`
}

type CommunityInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Locked bool   `json:"locked"`
}

func Validate(ctx context.Context, db *pgxpool.Pool, r *http.Request) (*PlayerValidationResponse, error) {
	user := ctx.Value(middleware.CtxUser).(middleware.UserContext)
	return &PlayerValidationResponse{
		Battletag: user.Battletag,
		IsAdmin:   user.CommunityRank <= user.CommunityOfficerRank,
		Community: CommunityInfo{
			Id:     user.CommunityId,
			Name:   user.CommunityName,
			Locked: user.CommunityLocked,
		},
	}, nil
}
