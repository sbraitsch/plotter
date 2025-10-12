package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
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

func (s *Server) AddPlayer(w http.ResponseWriter, r *http.Request) {
	req := &service.AddPlayerRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	player, err := service.AddPlayer(r.Context(), s.DB, req.Name)
	if err != nil {
		http.Error(w, "Failed to add player: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, player)
}

func (s *Server) InitializePlayerData(w http.ResponseWriter, r *http.Request) {
	req := &service.InitializePlayersRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	players, err := service.InitializePlayerData(r.Context(), s.DB, req.Names)
	if err != nil {
		http.Error(w, "Failed to initialize player data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, players)
}

func (s *Server) ListPlayerData(w http.ResponseWriter, r *http.Request) {
	data, err := service.ListPlayerData(r.Context(), s.DB)
	if err != nil {
		http.Error(w, "Failed to fetch player data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *Server) ListPlayerIds(w http.ResponseWriter, r *http.Request) {
	data, err := service.ListPlayerIds(r.Context(), s.DB)
	if err != nil {
		http.Error(w, "Failed to fetch player ids: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Failed to write response:", err)
	}
}
