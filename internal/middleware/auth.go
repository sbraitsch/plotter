package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserContext struct {
	Battletag            string
	CommunityId          string
	CommunityName        string
	CommunityRank        int
	CommunityOfficerRank int
	CommunityLocked      bool
	AccessToken          string
	Expiry               time.Time
}

type contextKey string

const (
	CtxUser contextKey = "user"
)

func TokenAuth(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Token")
			if token == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			var (
				battletag       string
				communityId     sql.NullString
				communityName   sql.NullString
				officerRank     sql.NullInt64
				communityLocked sql.NullBool
				communityRank   int
				accessToken     string
				expiry          time.Time
			)

			err := db.QueryRow(r.Context(),
				`SELECT
					u.battletag,
					u.community_id,
					c.name AS community_name,
					c.officer_rank,
					c.locked,
					u.community_rank,
					u.access_token,
					u.expiry
				FROM users u
				LEFT JOIN communities c
					ON u.community_id = c.id
				WHERE u.session_id = $1`,
				token,
			).Scan(&battletag, &communityId, &communityName, &officerRank, &communityLocked, &communityRank, &accessToken, &expiry)

			if err != nil {
				if err == pgx.ErrNoRows {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				log.Printf("Failed to find user matching given token: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			user := UserContext{
				Battletag:            battletag,
				CommunityId:          communityId.String,
				CommunityName:        communityName.String,
				CommunityRank:        communityRank,
				CommunityOfficerRank: int(officerRank.Int64),
				CommunityLocked:      communityLocked.Bool,
				AccessToken:          accessToken,
				Expiry:               expiry,
			}

			ctx := context.WithValue(r.Context(), CtxUser, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminAuth(db *pgxpool.Pool) func(http.Handler) http.Handler {
	tokenAuth := TokenAuth(db)
	return func(next http.Handler) http.Handler {
		return tokenAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(CtxUser).(UserContext)
			if !ok || (user.CommunityRank > user.CommunityOfficerRank) {
				http.Error(w, "forbidden: admin only", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}
