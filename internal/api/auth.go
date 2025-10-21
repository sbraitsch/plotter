package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sbraitsch/plotter/internal/service"
	"github.com/sbraitsch/plotter/internal/service/oauth"
	"github.com/sbraitsch/plotter/internal/storage"
	"golang.org/x/oauth2"
)

type AuthAPI interface {
	Routes(tmw func(http.Handler) http.Handler) chi.Router
}

type authAPIImpl struct {
	service  service.UserService
	oauthCfg *oauth2.Config
}

func NewAuthAPI(storage *storage.StorageClient, cfg *oauth2.Config) AuthAPI {
	return &authAPIImpl{service: service.NewUserService(storage), oauthCfg: cfg}
}

func (api *authAPIImpl) Routes(tmw func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/login", api.battleNetLogin)
	r.Get("/callback", api.battleNetCallback)
	r.With(tmw).Get("/guilds", api.listAvailableCommunities)

	return r
}

func (api *authAPIImpl) battleNetLogin(w http.ResponseWriter, r *http.Request) {
	url := api.oauthCfg.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusFound)
}

func (api *authAPIImpl) battleNetCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	sessionToken, err := api.service.RegisterUser(code, api.oauthCfg)

	if err != nil {
		http.Error(w, "Error creating new user.", http.StatusInternalServerError)
		return
	}

	frontendURL := fmt.Sprintf("%s/auth/success?token=%s", os.Getenv("FRONTEND_URL"), sessionToken)
	http.Redirect(w, r, frontendURL, http.StatusSeeOther)
}

func (api *authAPIImpl) listAvailableCommunities(w http.ResponseWriter, r *http.Request) {
	list, err := api.service.ListAvailableCommunities(r.Context())
	if err != nil {
		var tokenErr *oauth.TokenExpiredError
		if ok := errors.As(err, &tokenErr); ok {
			url := api.oauthCfg.AuthCodeURL("state")
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
	}

	if err != nil {
		http.Error(w, "Failed to get guild data", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, list)
}
