package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mathieuhays/uptime"
	"github.com/mathieuhays/uptime/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("missing env var: DATABASE_URL")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err = uptime.Migrate(db); err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing env var: JWT_SECRET")
	}

	apiConfig, err := uptime.NewApiConfig(database.New(db), jwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	const addr = "localhost:8080"
	router, err := uptime.NewRouter(apiConfig)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 5,
	}

	log.Printf("Starting server on %s", addr)
	log.Fatal(server.ListenAndServe())
}
