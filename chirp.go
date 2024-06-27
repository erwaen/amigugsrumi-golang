package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/erwaen/Chirpy/database"
)

type parameters struct {
	Body string `json:"body"`
}
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
			if err == database.ErrNotExist{
				respondWithError(w, http.StatusNotFound, "Chirp Not found")
			} else {
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error retrieving chirp: %s", err))
			}
			return
		}
		respondWithJson(w, http.StatusOK, chirp)
		return
	}
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps: %s", err))
		return
	}
	respondWithJson(w, 200, chirps)
}

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("%s", err))
		return
	}
	if params.Body == "" {
		respondWithError(w, 400, "no body field")
		return
	}
	if len(params.Body) >= 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	BLACK_WORDS := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	wordsList := strings.Split(params.Body, " ")
	for i, word := range wordsList {
		if _, exists := BLACK_WORDS[strings.ToLower(word)]; exists {
			wordsList[i] = "****"
		}
	}

	cleanedBody := strings.Join(wordsList, " ")

	// Save the chirp to the database
	newChirp, err := cfg.db.CreateChirp(cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirp: %s", err))
		return
	}

	respondWithJson(w, 201, newChirp)
}
