package main

import (
	"github.com/mathieuhays/uptime"
	"log"
	"net/http"
	"time"
)

func main() {
	const addr = "localhost:8080"
	router, err := uptime.NewRouter()
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: time.Minute,
	}

	log.Printf("Starting server on %s", addr)
	log.Fatal(server.ListenAndServe())
}
