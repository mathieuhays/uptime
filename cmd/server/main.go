package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mathieuhays/uptime"
	"github.com/mathieuhays/uptime/internal/database"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func run(getenv func(string) string, stdout, stderr io.Writer) error {
	dbURL := getenv("DATABASE_URL")
	if dbURL == "" {
		return errors.New("missing env var: DATABASE_URL")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	if err = uptime.Migrate(db); err != nil {
		return err
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("missing env var: JWT_SECRET")
	}

	dbQueries := database.New(db)

	apiConfig, err := uptime.NewApiConfig(jwtSecret)
	if err != nil {
		return err
	}

	logger := log.New(stderr, "", log.LstdFlags|log.Lshortfile)

	const addr = "localhost:8080"
	router, err := uptime.NewRouter(logger, dbQueries, apiConfig)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 5,
	}

	_, _ = fmt.Fprintf(stdout, "Starting server on %s\n", addr)
	if err = server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if err := run(os.Getenv, os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
