package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sbraitsch/plotter/internal/service"
)

type InitializePlayersRequest struct {
	Players []service.PlayerInit `json:"players"`
}

type AddPlayerRequest struct {
	Name string `json:"name"`
}

type PlotMappingRequest struct {
	Name     string      `json:"name"`
	PlotData map[int]int `json:"plotData"`
}

func (s *Server) Validate(w http.ResponseWriter, r *http.Request) {
	player, err := service.Validate(r.Context(), s.DB, r)
	if err != nil {
		http.Error(w, "Failed to add player: "+err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, player)
}

func (s *Server) UpdateMapping(w http.ResponseWriter, r *http.Request) {
	req := &PlotMappingRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	player, err := service.Validate(r.Context(), s.DB, r)

	if err != nil {
		http.Error(w, "Token Validation Failed", http.StatusUnauthorized)
		return
	}

	err = service.ValidateNameToken(player.Name, req.Name)

	if err != nil {
		http.Error(w, "Token and Player do not match: "+err.Error(), http.StatusUnauthorized)
		return
	}

	err = service.UpdatePlayerData(r.Context(), s.DB, req.Name, req.PlotData)

	if err != nil {
		http.Error(w, "Failed to update player data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.ListPlayerData(w, r)
}

func (s *Server) AddPlayer(w http.ResponseWriter, r *http.Request) {
	req := &AddPlayerRequest{}

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
	req := &InitializePlayersRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	players, err := service.InitializePlayerData(r.Context(), s.DB, req.Players)
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
