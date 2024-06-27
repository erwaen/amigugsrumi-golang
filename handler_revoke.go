package main

import (
	"errors"
	"net/http"

	"github.com/erwaen/Chirpy/auth"
	"github.com/erwaen/Chirpy/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	_, err = cfg.db.RevokeRefreshToken(refreshToken)
    if err!= nil {
        if errors.Is(err, database.ErrNotExist){
            respondWithError(w, http.StatusUnauthorized, "Refresh Token doesn't exist")
        }else{
			respondWithError(w, http.StatusUnauthorized, "Couldn't revoke the token")
        }
        return 
    }
    respondWithoutJson(w, http.StatusNoContent)
}
