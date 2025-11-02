package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sbraitsch/plotter/internal/model"
	"github.com/sbraitsch/plotter/internal/service"
	"github.com/sbraitsch/plotter/internal/storage"
)

type UserAPI interface {
	Routes() chi.Router
}

type userAPIImpl struct {
	service service.UserService
}

func NewUserAPI(storage *storage.StorageClient) UserAPI {
	return &userAPIImpl{service: service.NewUserService(storage)}
}

func (api *userAPIImpl) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/validate", api.validate)
	r.Post("/update", api.updatePlayerData)

	return r
}

func (api *userAPIImpl) validate(w http.ResponseWriter, r *http.Request) {
	player, err := api.service.Validate(r.Context())
	if err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, "Failed to validate session", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, player)
}

func (api *userAPIImpl) updatePlayerData(w http.ResponseWriter, r *http.Request) {
	req := &model.PlayerUpdateRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := api.service.SetNote(r.Context(), req.Note); err != nil {
		http.Error(w, "Failed to update player note", http.StatusInternalServerError)
		return
	}

	updated, err := api.service.UpdateMappings(r.Context(), req.PlotData)

	if err != nil {
		http.Error(w, "Failed to update player data", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, updated)
}
