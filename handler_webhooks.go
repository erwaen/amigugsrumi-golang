package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/erwaen/Chirpy/auth"
	"github.com/erwaen/Chirpy/database"
)

func (cfg *apiConfig) handlerWBUpgrade(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID int `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	token, err := auth.GetApiKeyToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token")
		return
	}
	if token != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Not allowed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithoutJson(w, 204)
		return
	}

	userID := params.Data.UserID
	_, err = cfg.db.UpgradeUserRed(userID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, 404, "User not found")
		} else {
			respondWithError(w, 500, "Couldn't retrieve user")
		}
		return
	}
	respondWithoutJson(w, 204)

}
