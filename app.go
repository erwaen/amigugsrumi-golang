package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/erwaen/Chirpy/tursodb"
	"github.com/erwaen/Chirpy/types"

	"github.com/erwaen/Chirpy/database"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwtSecret      string
	polkaKey       string
	tursoDB        *tursodb.TursoDB
}

func main() {
	// Initialize the databasenpblock
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if jwtSecret == "" {
		log.Fatal("polka key environment variable is not set")
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	tursoUrl := os.Getenv("TURSO_DATABASE_URL")
	if tursoUrl == "" {
		log.Fatal("TURSO_DATABASE_URL key environment variable is not set")
	}
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
	if tursoToken == "" {
		log.Fatal("TURSO_AUTH_TOKEN key environment variable is not set")
	}
	completeUrl := tursoUrl + "?authToken=" + tursoToken
	tursoDB, err := sql.Open("libsql", completeUrl)
	if err != nil {
		log.Fatal("error in connect turso db", err)
	}
	defer tursoDB.Close()
	tursoDBWrapper := tursodb.NewTursoDB(tursoDB)

	apiCfg := apiConfig{
		fileserverHits: 0,
		db:             db,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
		tursoDB:        tursoDBWrapper,
	}
	mux := http.NewServeMux()
	fhandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.Handle("/app/*", fhandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerNewChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerReadChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerReadChirps)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.handlerDeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerWBUpgrade)

	mux.HandleFunc("POST /api/users", apiCfg.handlerNewUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("GET /api/tursousers", apiCfg.handlerTursoUsers)
	mux.HandleFunc("GET /api/tursoitems", apiCfg.handlerTursoItems)
	mux.HandleFunc("GET /api/tursoitemsstock", apiCfg.handlerTursoItemsStock)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", ".", "8080")
	log.Fatal(server.ListenAndServe())

}

type TursoUser struct {
	ID   int
	Name string
}

func (cfg *apiConfig) handlerTursoUsers(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Users []types.TursoUser `json:"users"`
	}

	users, err := cfg.tursoDB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting users: %s", err))
		return
	}

	respondWithJson(w, http.StatusOK, response{
		Users: users,
	})

}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (cfg *apiConfig) handlerTursoItems(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Items []types.TursoItem `json:"items"`
	}

	items, err := cfg.tursoDB.GetItems()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting items: %s", err))
		return
	}

	respondWithJson(w, http.StatusOK, response{
		Items: items,
	})

}

func (cfg *apiConfig) handlerTursoItemsStock(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	type response struct {
		Items []types.TursoItemStock `json:"items"`
	}
	items, err := cfg.tursoDB.GetItemsStock()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting items: %s", err))
		return
	}
	respondWithJson(w, http.StatusOK, response{
		Items: items,
	})
}
