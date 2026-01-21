package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{}

	smux := http.NewServeMux()

	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	smux.Handle("/app/", fileServer)
	smux.HandleFunc("GET /api/healthz", handlerReadiness)
	smux.HandleFunc("GET /api/metrics", apiCfg.handlerHits)
	smux.HandleFunc("POST /api/reset", apiCfg.handlerReset)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: smux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
