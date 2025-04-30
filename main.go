package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Khazz0r/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading environment variables: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWTSECRET")
	if jwtSecret == "" {
		log.Fatal("JWT secret must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening chirpy database: %v", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("/assets", http.FileServer(http.Dir("logo.png")))

	mux.HandleFunc("GET /api/healthz", handlerOkStatus)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerDeleteAllUsers)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerNumOfRequests)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
