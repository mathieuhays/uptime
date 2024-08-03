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
	"net"
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

	// should probably be configurable via env vars
	const host = "localhost"
	const port = "8080"

	userStore := uptime.NewUserStore(dbQueries)
	sessionStore := uptime.NewSessionStore(dbQueries, apiConfig)

	srv := uptime.NewServer(logger, userStore, sessionStore, apiConfig)

	server := &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           srv,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 5,
	}

	_, _ = fmt.Fprintf(stdout, "Starting server on %s\n", server.Addr)
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
