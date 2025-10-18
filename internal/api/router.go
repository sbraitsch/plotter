package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"

	"github.com/sbraitsch/plotter/internal/middleware"
)

type Config struct {
	DbUrl        string
	Port         string
	ClientId     string
	ClientSecret string
}

type Server struct {
	DB    *pgxpool.Pool
	Oauth *oauth2.Config
}

func NewServer(db *pgxpool.Pool, cfg Config) Server {
	bnetOAuthConfig := &oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://eu.battle.net/oauth/authorize",
			TokenURL: "https://eu.battle.net/oauth/token",
		},
		Scopes: []string{"wow.profile"},
	}

	return Server{DB: db, Oauth: bnetOAuthConfig}

}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://plotter.sbraitsch.dev", "http://localhost:3000"}, // Production
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.LoggingMiddleware)

	//oauth
	r.Get("/auth/bnet/login", s.BattleNetLogin)
	r.Get("/auth/bnet/callback", s.BattleNetCallback)

	// user endpoints
	r.Group(func(user chi.Router) {
		user.Use(middleware.TokenAuth(s.DB))
		user.Get("/community", s.GetCommunity)
		user.Get("/assignments", s.GetAssignments)
		user.Post("/join", s.JoinCommunity)
		user.Get("/guilds", s.GetUserGuilds)
		user.Get("/validate", s.Validate)
		user.Post("/update", s.UpdateMapping)
	})

	// admin-only endpoints
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.AdminAuth(s.DB))
		admin.Get("/optimize", s.RunOptimizer)
		admin.Post("/lock", s.LockAssignments)
		admin.Post("/config", s.SetOfficerRank)
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}
