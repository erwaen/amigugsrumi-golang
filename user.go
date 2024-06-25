package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/erwaen/Chirpy/database"
	"github.com/erwaen/Chirpy/types"
	"golang.org/x/crypto/bcrypt"
)

type parameterNewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameterNewUser{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	if params.Email == "" || params.Password == "" {
		respondWithError(w, 400, "no email or password field")
		return
	}

	// Save the user to the database
	newUser, err := cfg.db.CreateUser(params.Email, params.Password)
	if err != nil {
		if err == database.ErrUserAlreadyExist {

			respondWithError(w, 400, fmt.Sprintf("Error creating user: %s", err))
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %s", err))
		}
		return
	}

	respondWithJson(w, 201, newUser)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameterNewUser{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	if params.Email == "" || params.Password == "" {
		respondWithError(w, 400, "no email or password field")
		return
	}

	// get the user from the database
	user, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		if err == database.ErrUserAlreadyExist {
			respondWithError(w, 400, fmt.Sprintf("Error logging user: %s", err))
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error logging user: %s", err))
		}
		return
	}

	//compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Error logging user: %s", err))
	}
    loggedUser:= types.LoggedUser{
        Id: user.Id,
        Email: user.Email,
    }

	respondWithJson(w, 200, loggedUser)
}
