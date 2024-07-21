package main

import (
	"github.com/mathieuhays/uptime"
	"log"
	"net/http"
)

func main() {
	const addr = "localhost:8080"
	server, err := uptime.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
