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
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to sql database: %s", err)
	}

	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{dbQueries: database.New(db), platform: platform}

	smux := http.NewServeMux()

	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	smux.Handle("/app/", fileServer)

	smux.HandleFunc("GET /api/healthz", handlerReadiness)

	smux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	smux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	smux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	smux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	smux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: smux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
