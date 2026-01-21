package main

import (
	"fmt"
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

	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	smux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServer))
	smux.HandleFunc("/healthz", handlerReadiness)
	smux.HandleFunc("/metrics", apiCfg.handlerHits)
	smux.HandleFunc("/reset", apiCfg.handlerHitsReset)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: smux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (apiCfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf("Hits: %d", apiCfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (apiCfg *apiConfig) handlerHitsReset(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileserverHits.Store(0)
}
