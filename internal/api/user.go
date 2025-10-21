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

	r.Get("/{id}", api.getUserByToken)
	r.Get("/validate", api.validate)
	r.Post("/update", api.updateMapping)

	return r
}

func (api *userAPIImpl) getUserByToken(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "id")
	user, err := api.service.GetUserByToken(r.Context(), userId)
	if err != nil {
		http.Error(w, "Failed to retrieve user associated with token", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
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

func (api *userAPIImpl) updateMapping(w http.ResponseWriter, r *http.Request) {
	req := &model.PlotMappingRequest{}

	if err := render.Decode(r, req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := api.service.UpdateMappings(r.Context(), req.PlotData)

	if err != nil {
		http.Error(w, "Failed to update player data", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, updated)
}
