package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/erwaen/Chirpy/auth"
	"github.com/erwaen/Chirpy/database"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	rfStruct, err := cfg.db.GetRefreshTokenStruct(refreshToken)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusUnauthorized, "Refresh token doesn't exist")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token from db")
		}
		return
	}
	if time.Now().After(rfStruct.ExpireAt) {
		respondWithError(w, http.StatusUnauthorized, "Refresh Token expired")
		return
	}

	userID := rfStruct.UserID
	user, err := cfg.db.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusUnauthorized, "User doesnt exist")
		} else {
			respondWithError(w, http.StatusUnauthorized, "Couldn't Retrieve de user")
		}
		return
	}

	defaultExpiration := 60 * 60
	token, err := auth.MakeJWT(user.Id, cfg.jwtSecret, time.Duration(defaultExpiration)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJson(w, 200, response{
		Token: token,
	})

}
