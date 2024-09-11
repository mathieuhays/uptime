package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mathieuhays/uptime"
	"github.com/mathieuhays/uptime/internal/healthcheck"
	"github.com/mathieuhays/uptime/internal/website"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/tursodatabase/go-libsql"
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

	db, err := sql.Open("libsql", "file:"+databasePath)
	if err != nil {
		log.Fatalf("database error: %s", err)
	}

	if err = uptime.Migrate(db, "turso"); err != nil {
		log.Fatalf("migration error: %s", err)
	}

	websiteRepository, err := website.NewSQLiteRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	healthCheckRepo, err := healthcheck.NewSQLiteRepository(db)
	if err != nil {
		log.Fatal(err)
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

	crawler := uptime.NewCrawler(
		healthCheckRepo,
		websiteRepository,
		time.Minute,
		5,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go crawler.Start(ctx)

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
