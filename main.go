package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{}

	smux := http.NewServeMux()

	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	smux.Handle("/app/", fileServer)

	smux.HandleFunc("GET /api/healthz", handlerReadiness)
	smux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	smux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	smux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: smux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
