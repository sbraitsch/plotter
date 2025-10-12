package api

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB *pgxpool.Pool
}

func (s *Server) Start(addr string) error {

	http.HandleFunc("/update", s.Update)
	http.HandleFunc("/pull", s.Pull)

	return http.ListenAndServe(addr, nil)
}
