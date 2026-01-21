package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf("Hits: %d", apiCfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
