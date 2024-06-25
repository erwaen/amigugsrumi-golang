package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}
type returnError struct {
	Error string `json:"error"`
}
type returnValid struct {
	CleanedBody string `json:"cleaned_body"`
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

func handlerNewChirp(w http.ResponseWriter, r *http.Request) {
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

	valid := returnValid{
		CleanedBody: strings.Join(wordsList, " "),
	}

	respondWithJson(w, 200, valid)
}
