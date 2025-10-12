package api

import (
	"encoding/json"
	"net/http"

	"github.com/sbraitsch/plotter/internal/service"
)

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "update successful"})
}

func (s *Server) Pull(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": "pull result"})

}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	service.CreatePlayer(s.DB, w, r)
}

func (s *Server) ListPlayers(w http.ResponseWriter, r *http.Request) {
	service.ListPlayers(s.DB, w, r)
}
