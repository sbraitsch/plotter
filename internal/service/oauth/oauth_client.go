package oauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/sbraitsch/plotter/internal/middleware"
	"golang.org/x/oauth2"
)

type TokenExpiredError struct{}

func (e *TokenExpiredError) Error() string {
	return "token expired, user must re-authenticate"
}

func tokenFetcher(ctx context.Context) (*oauth2.Token, error) {
	user, ok := ctx.Value(middleware.CtxUser).(middleware.UserContext)

	if !ok {
		return nil, errors.New("Invalid user context")
	}

	return &oauth2.Token{
		AccessToken: user.AccessToken,
		Expiry:      user.Expiry,
	}, nil
}

func refreshToken() (*oauth2.Token, error) {
	return nil, &TokenExpiredError{}
}

func GetClient(ctx context.Context) *http.Client {
	ts := &BattleNetTokenSource{
		TokenFunc: func() (*oauth2.Token, error) {
			return tokenFetcher(ctx)
		},
		RefreshFunc: refreshToken,
	}

	return oauth2.NewClient(ctx, ts)
}
