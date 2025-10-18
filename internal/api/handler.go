package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sbraitsch/plotter/internal/middleware"
	"github.com/sbraitsch/plotter/internal/service"
	"github.com/sbraitsch/plotter/internal/service/oauth"
)

type PlotMappingRequest struct {
	PlotData map[int]int `json:"plotData"`
}

type CommunityRankRequest struct {
	MinRank int `json:"minRank"`
}

func (s *Server) Validate(w http.ResponseWriter, r *http.Request) {
	player, err := service.Validate(r.Context(), s.DB, r)
	if err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, "Failed to validate session", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, player)
}

func (s *Server) UpdateMapping(w http.ResponseWriter, r *http.Request) {
	req := &PlotMappingRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := service.UpdatePlayerData(r.Context(), s.DB, req.PlotData)

	if err != nil {
		http.Error(w, "Failed to update player data", http.StatusInternalServerError)
		return
	}

	s.GetCommunity(w, r)
}

func (s *Server) GetCommunity(w http.ResponseWriter, r *http.Request) {
	data, err := service.GetCommunity(r.Context(), s.DB)
	if err != nil {
		http.Error(w, "Failed to get community data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *Server) JoinCommunity(w http.ResponseWriter, r *http.Request) {
	communityId := r.URL.Query().Get("community")
	if communityId == "" {
		http.Error(w, "Missing community parameter", http.StatusBadRequest)
		return
	}
	err := service.JoinCommunity(r.Context(), s.DB, communityId)
	if err != nil {
		http.Error(w, "Failed to join community", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) RunOptimizer(w http.ResponseWriter, r *http.Request) {
	community, err := service.GetCommunity(r.Context(), s.DB)
	if err != nil {
		http.Error(w, "Failed to run optimizer", http.StatusInternalServerError)
		return
	}
	optimized := service.Optimize(community)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(optimized); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *Server) LockAssignments(w http.ResponseWriter, r *http.Request) {
	// unlock if locked
	user := r.Context().Value(middleware.CtxUser).(middleware.UserContext)
	if user.CommunityLocked {
		err := service.UnlockCommunity(r.Context(), s.DB, user.CommunityId)
		if err != nil {
			http.Error(w, "Failed to unlock community", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	community, err := service.GetCommunity(r.Context(), s.DB)
	if err != nil {
		http.Error(w, "Failed to run optimizer", http.StatusInternalServerError)
		return
	}

	assignments := service.Optimize(community)

	err = service.SaveAssignmentsAndLock(r.Context(), s.DB, assignments, community.Id)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(assignments); err != nil {
		log.Println("Failed to write response:", err)
	}
}
func (s *Server) GetAssignments(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.CtxUser).(middleware.UserContext)
	data, err := service.GetAssignments(r.Context(), s.DB, user.CommunityId)

	if err != nil {
		http.Error(w, "Failed to get plot assignments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *Server) GetUserGuilds(w http.ResponseWriter, r *http.Request) {
	data, err := service.GetGuilds(r.Context(), s.DB)
	if err != nil {
		var tokenErr *oauth.TokenExpiredError
		if ok := errors.As(err, &tokenErr); ok {
			url := s.Oauth.AuthCodeURL("state")
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
	}

	if err != nil {
		http.Error(w, "Failed to get guild data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *Server) SetOfficerRank(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.CtxUser).(middleware.UserContext)
	req := &CommunityRankRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MinRank < 0 {
		http.Error(w, "Nice try.", http.StatusBadRequest)
		return
	}

	_, err := s.DB.Exec(r.Context(),
		`UPDATE communities
		 SET officer_rank = $1
		 WHERE id = $2`,
		req.MinRank, user.CommunityId,
	)

	if err != nil {
		http.Error(w, "Failed to update community setting", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
