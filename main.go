package main

import (
	"net/http"
	"sync/atomic"
	"fmt"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	var apiCfg apiConfig

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("/assets", http.FileServer(http.Dir("logo.png")))
	mux.HandleFunc("GET /healthz", handlerOkStatus)
	mux.HandleFunc("GET /metrics", apiCfg.handlerNumOfRequests)
	mux.HandleFunc("POST /reset", apiCfg.handlerResetNumOfRequests)

	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	server.ListenAndServe()
}

// middleware to track how many times the website has been accessed
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

// handler to show an ok status when /healthz is accessed
func handlerOkStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	message := "OK"

	w.Write([]byte(message))
}

// handler to write out the number of hits that have happened with when /metrics is accessed
func (cfg *apiConfig) handlerNumOfRequests(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	hits := cfg.fileserverHits.Load()
	message := fmt.Sprintf("Hits: %d", hits)

	w.Write([]byte(message))
}

// handler to reset number of hits back to 0 when /reset is accessed
func (cfg *apiConfig) handlerResetNumOfRequests(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)

	w.Write([]byte("Hits have been reset"))
}
