package oauth

import (
	"errors"
	"time"

	"golang.org/x/oauth2"
)

type BattleNetTokenSource struct {
	TokenFunc   func() (*oauth2.Token, error) // fetches a fresh token (e.g., from DB or memory)
	RefreshFunc func() (*oauth2.Token, error) // called when token is expired
}

func (b *BattleNetTokenSource) Token() (*oauth2.Token, error) {
	token, err := b.TokenFunc()
	if err != nil {
		return nil, err
	}

	if token == nil || token.Expiry.IsZero() || time.Now().After(token.Expiry) {
		if b.RefreshFunc != nil {
			return b.RefreshFunc()
		}
		return nil, errors.New("access token expired, re-authentication required")
	}

	return token, nil
}
