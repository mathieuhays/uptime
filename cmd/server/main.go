package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mathieuhays/uptime"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	errMissingHostname = errors.New("missing HOSTNAME env variable")
	errMissingPort     = errors.New("missing PORT env variable")
)

func run(getenv func(string) string, stdout, stderr io.Writer) error {
	hostname := getenv("HOSTNAME")
	port := getenv("PORT")

	if hostname == "" {
		return errMissingHostname
	}

	if port == "" {
		return errMissingPort
	}

	templ, err := uptime.TemplateEngine()
	if err != nil {
		log.Fatalf("error loading templates: %s", err)
	}

	serverHandler := uptime.NewServer(templ)

	server := &http.Server{
		Addr:              net.JoinHostPort(hostname, port),
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
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
