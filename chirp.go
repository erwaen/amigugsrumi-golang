package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/erwaen/Chirpy/auth"
	"github.com/erwaen/Chirpy/database"
)

type returnError struct {
	Error string `json:"error"`
}
type returnValid struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	log.Printf(msg)
	error := returnError{
		Error: msg,
	}
	respondWithJson(w, code, error)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithoutJson(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func (cfg *apiConfig) handlerReadChirps(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString != "" {
		id, err := strconv.Atoi(idString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid ID parameter")
			return
		}

		chirp, err := cfg.db.GetChirp(id)
		if err != nil {
			if err == database.ErrNotExist {
				respondWithError(w, http.StatusNotFound, "Chirp Not found")
			} else {
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error retrieving chirp: %s", err))
			}
			return
		}
		respondWithJson(w, http.StatusOK, chirp)
		return
	}

	s := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	authorID, err := strconv.Atoi(s)
	if err != nil {
		authorID = 0
	}

	chirps, err := cfg.db.GetChirps(authorID, sort)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps: %s", err))
		return
	}
	respondWithJson(w, 200, chirps)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	idString := r.PathValue("id")
	chirpID, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	chirp, err := cfg.db.GetChirp(chirpID)
	if err != nil {
		if err == database.ErrNotExist {
			respondWithError(w, http.StatusNotFound, "Chirp Not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get the chirp")
		}
		return
	}

	if chirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "You are not allowed to delete this chirp")
		return
	}
	_, err = cfg.db.DeleteChirp(chirpID)
	if err != nil {
		if err == database.ErrNotExist {
			respondWithError(w, http.StatusNotFound, "Chirp Not found when trying to delete")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't delete the chirp")
		}
	}
	respondWithoutJson(w, http.StatusNoContent)
}

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		ID       int    `json:"id"`
		Body     string `json:"body"`
		AuthorID int    `json:"author_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Save the chirp to the database
	newChirp, err := cfg.db.CreateChirp(cleaned, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}
	respondWithJson(w, http.StatusCreated, response{
		ID:       newChirp.Id,
		Body:     newChirp.Body,
		AuthorID: newChirp.AuthorID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	BLACK_WORDS := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, BLACK_WORDS)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
