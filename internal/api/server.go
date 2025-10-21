package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"

	"github.com/sbraitsch/plotter/internal/middleware"
	"github.com/sbraitsch/plotter/internal/storage"
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

	storageClient := storage.NewStorageClient(s.DB)
	tokenMiddleware := middleware.TokenAuth(storageClient)
	adminMiddleware := middleware.AdminAuth(storageClient)
	userAPI := NewUserAPI(storageClient)
	communityAPI := NewCommunityAPI(storageClient)
	authApi := NewAuthAPI(storageClient, s.Oauth)

	r.Route("/user", func(r chi.Router) {
		r.Use(tokenMiddleware)
		r.Mount("/", userAPI.Routes())
	})

	r.Route("/community", func(r chi.Router) {
		r.Mount("/", communityAPI.Routes(tokenMiddleware, adminMiddleware))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Mount("/bnet", authApi.Routes(tokenMiddleware))
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}
