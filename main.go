package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", handlerReadiness)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())

}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
