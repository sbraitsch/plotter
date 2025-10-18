package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func (s *Server) BattleNetLogin(w http.ResponseWriter, r *http.Request) {
	url := s.Oauth.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Server) BattleNetCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	ctx := context.Background()

	token, err := s.Oauth.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := s.Oauth.Client(ctx, token)
	resp, err := client.Get("https://oauth.battle.net/oauth/userinfo")
	if err != nil {
		http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var profile struct {
		Battletag string `json:"battletag"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		http.Error(w, "Failed to decode profile", http.StatusInternalServerError)
		return
	}

	sessionToken := uuid.New().String()

	_, err = s.DB.Exec(r.Context(), `INSERT INTO users(battletag, session_id, access_token, expiry)
                      VALUES($1, $2, $3, $4)
                      ON CONFLICT(battletag) DO UPDATE
                      SET session_id=$2, access_token=$3, expiry=$4`,
		profile.Battletag, sessionToken, token.AccessToken, token.Expiry)

	if err != nil {
		log.Printf("Failed to insert new user: %v", err)
		http.Error(w, "Error inserting new user.", http.StatusInternalServerError)
		return
	}

	frontendURL := fmt.Sprintf("%s/auth/success?token=%s", os.Getenv("FRONTEND_URL"), sessionToken)
	http.Redirect(w, r, frontendURL, http.StatusSeeOther)
}
