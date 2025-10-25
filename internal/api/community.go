package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sbraitsch/plotter/internal/middleware"
	"github.com/sbraitsch/plotter/internal/model"
	"github.com/sbraitsch/plotter/internal/service"
	"github.com/sbraitsch/plotter/internal/storage"
)

type CommunityAPI interface {
	Routes(tmw, amw func(http.Handler) http.Handler) chi.Router
}

type communityAPIImpl struct {
	service service.CommunityService
}

func NewCommunityAPI(storage *storage.StorageClient) CommunityAPI {
	return &communityAPIImpl{service: service.NewCommunityService(storage)}
}

func (api *communityAPIImpl) Routes(tmw, amw func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Group(func(user chi.Router) {
		user.Use(tmw)
		user.Get("/", api.getCommunityData)
		user.Post("/join/{id}", api.joinCommunity)
	})

	r.Group(func(admin chi.Router) {
		admin.Use(amw)
		admin.Get("/optimize", api.runOptimizer)
		admin.Get("/assignments", api.getAssignments)
		admin.Post("/lock", api.toggleCommunityLock)
		admin.Post("/config", api.setOfficerRank)
	})

	return r
}

func (api *communityAPIImpl) getCommunityData(w http.ResponseWriter, r *http.Request) {
	community, err := api.service.GetCommunityData(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve community data", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, community)
}

func (api *communityAPIImpl) joinCommunity(w http.ResponseWriter, r *http.Request) {
	communityId := chi.URLParam(r, "id")
	err := api.service.JoinCommunity(r.Context(), communityId)
	if err != nil {
		http.Error(w, "Failed to join community: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *communityAPIImpl) runOptimizer(w http.ResponseWriter, r *http.Request) {
	community, err := api.service.GetCommunityData(r.Context())
	if err != nil {
		http.Error(w, "Failed to run optimizer", http.StatusInternalServerError)
		return
	}
	optimized := community.Optimize()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(optimized); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (api *communityAPIImpl) toggleCommunityLock(w http.ResponseWriter, r *http.Request) {
	// unlock if locked
	user := r.Context().Value(middleware.CtxUser).(*model.User)
	assignments, err := api.service.ToggleCommunityLock(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to lock community.", http.StatusInternalServerError)
	}
	if assignments != nil {
		render.JSON(w, r, assignments)
	}

	w.WriteHeader(http.StatusOK)
}

func (api *communityAPIImpl) getAssignments(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.CtxUser).(*model.User)
	assignments, err := api.service.GetAssignments(r.Context(), user.Community.Id)

	if err != nil {
		http.Error(w, "Failed to get plot assignments", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, assignments)
}

func (api *communityAPIImpl) setOfficerRank(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.CtxUser).(*model.User)
	req := &model.CommunityRankRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MinRank < 0 {
		http.Error(w, "Nice try.", http.StatusBadRequest)
		return
	}

	err := api.service.SetOfficerRank(r.Context(), user.Community.Id, req.MinRank)

	if err != nil {
		log.Printf("Failed to update community settings: %v", err)
		http.Error(w, "Error updating community settings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
