package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/erwaen/Chirpy/auth"
)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	// get the user from the database
	user, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

    defaultExpiration := 60*60
    if params.ExpiresInSeconds ==0{
        params.ExpiresInSeconds = defaultExpiration
    } else if params.ExpiresInSeconds > defaultExpiration {
        params.ExpiresInSeconds = defaultExpiration
    }

    token, err:= auth.MakeJWT(user.Id, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
        return
    }
    refreshToken, err:= auth.MakeRefreshT()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create Refresh Token")
        return
    }

	respondWithJson(w, 200, response {
        User: User{
            ID: user.Id,
            Email: user.Email,
        },
        Token: token,
    })
}
