package api

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB *pgxpool.Pool
}

func (s *Server) Start(addr string) error {

	http.HandleFunc("/update", s.Update)
	http.HandleFunc("/pull", s.Pull)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Unmatched route: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

	return http.ListenAndServe(addr, nil)
}
