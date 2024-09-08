package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mathieuhays/uptime"
	"github.com/mathieuhays/uptime/internal/website"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	errMissingPort         = errors.New("missing PORT env variable")
	errMissingDatabasePath = errors.New("missing DATABASE_PATH env variable")
)

func run(getenv func(string) string, stdout, stderr io.Writer) error {
	port := getenv("PORT")
	databasePath := getenv("DATABASE_PATH")

	if port == "" {
		return errMissingPort
	}

	if databasePath == "" {
		return errMissingDatabasePath
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("database error: %s", err)
	}

	websiteRepository := website.NewSQLiteRepository(db)

	if err = websiteRepository.Migrate(); err != nil {
		log.Fatalf("website repo migration error: %s", err)
	}

	templ, err := uptime.TemplateEngine()
	if err != nil {
		log.Fatalf("error loading templates: %s", err)
	}

	serverHandler := uptime.NewServer(templ, websiteRepository)

	server := &http.Server{
		Addr:              net.JoinHostPort("", port),
		Handler:           serverHandler,
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
		log.Printf("env file error: %s", err)
	}

	if err := run(os.Getenv, os.Stdout, os.Stderr); err != nil {
		log.Printf("listen err: %s\n", err)
		os.Exit(1)
	}

	log.Println("goodbye")
}
