package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zombfeed/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file loaded: %w\n", err)
	}
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("API_SECRET_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to sql database: %s", err)
	}

	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{dbQueries: database.New(db), platform: platform, secret: secret}
	smux := http.NewServeMux()

	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	smux.Handle("/app/", fileServer)

	smux.HandleFunc("GET /api/healthz", handlerReadiness)

	smux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	smux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	smux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	smux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	smux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	smux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	smux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirpByID)

	smux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	smux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	smux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	smux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: smux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
