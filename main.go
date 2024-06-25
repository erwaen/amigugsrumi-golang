package main

import (
	"log"
	"net/http"

	"github.com/erwaen/Chirpy/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	// Initialize the database
	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		db:             db,
	}
	mux := http.NewServeMux()
	fhandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.Handle("/app/*", fhandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerNewChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerReadChirps)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
	log.Fatal(server.ListenAndServe())

}
