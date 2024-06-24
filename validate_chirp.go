package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}
type returnError struct {
	Error string `json:"error"`
}
type returnValid struct {
	Valid bool `json:"valid"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	error := returnError{}
	err := decoder.Decode(&params)
	if err != nil {
		error.Error = "Something went wrong"
		data, err := json.Marshal(error)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		w.Write(data)
		return
	}
	if len(params.Body) >= 140 {
		error.Error = "Chirp is to long"
		data, err := json.Marshal(error)
		if err != nil {
			log.Printf("Error marshalling JSON: %s\n", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Printf("chirp is to long %s\n", err)
		w.WriteHeader(400)
		w.Write(data)
		return
	}

    if params.Body == ""{
        
		w.WriteHeader(400)
		w.Write([]byte ("no body field"))
        return 
    }
	valid := returnValid{
		Valid: true,
	}
	data, err := json.Marshal(valid)
	if err != nil {
		log.Printf("Error marshalling valid response JSON: %s\n", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
