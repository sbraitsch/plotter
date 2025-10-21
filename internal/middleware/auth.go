package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/sbraitsch/plotter/internal/model"
	"github.com/sbraitsch/plotter/internal/storage"
)

type contextKey string

const (
	CtxUser contextKey = "user"
)

func TokenAuth(storage *storage.StorageClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Token")
			if token == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			user, err := storage.GetUserByToken(r.Context(), token)

			if err != nil {
				if err == pgx.ErrNoRows {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				log.Printf("Failed to find user matching given token: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUser, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminAuth(storage *storage.StorageClient) func(http.Handler) http.Handler {
	tokenAuth := TokenAuth(storage)
	return func(next http.Handler) http.Handler {
		return tokenAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(CtxUser).(*model.User)
			if !ok || (user.CommunityRank > user.Community.OfficerRank) {
				http.Error(w, "forbidden: admin only", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}
