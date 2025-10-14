package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidToken = errors.New("invalid token")

type PlayerValidationResponse struct {
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}

func Validate(ctx context.Context, db *pgxpool.Pool, r *http.Request) (*PlayerValidationResponse, error) {
	token := r.Header.Get("X-Token")
	if token == "" {
		return nil, errors.New("missing token header")
	}

	if _, err := uuid.Parse(token); err != nil {
		return nil, ErrInvalidToken
	}

	var name string
	var isAdmin bool
	err := db.QueryRow(ctx,
		`SELECT name, is_admin FROM players WHERE uuid = $1`,
		token,
	).Scan(&name, &isAdmin)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	return &PlayerValidationResponse{
		Name:    name,
		IsAdmin: isAdmin,
	}, nil
}

func ValidateNameToken(saved, requested string) error {
	if requested == "" || saved != requested {
		return ErrInvalidToken
	}
	return nil
}
