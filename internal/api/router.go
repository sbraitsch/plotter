package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB        *pgxpool.Pool
	AdminUUID string
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	// Public endpoints
	r.Post("/update", s.Update)
	r.Get("/pull", s.Pull)

	// Admin-only endpoints
	r.Group(func(admin chi.Router) {
		admin.Use(AdminAuth(s.AdminUUID))
		admin.Post("/create", s.Create)
		admin.Get("/players", s.ListPlayers)
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}
