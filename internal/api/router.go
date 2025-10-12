package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB        *pgxpool.Pool
	AdminUUID string
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://plotter.sbraitsch.dev", "http://localhost:3000"}, // Production
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Public endpoints
	r.Post("/update", s.Update)
	r.Get("/pull", s.Pull)
	r.Get("/playerdata", s.ListPlayerData)

	// Admin-only endpoints
	r.Group(func(admin chi.Router) {
		admin.Use(AdminAuth(s.AdminUUID))
		admin.Post("/create", s.AddPlayer)
		admin.Get("/playerids", s.ListPlayerIds)
		admin.Post("/initialize", s.InitializePlayerData)
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}
