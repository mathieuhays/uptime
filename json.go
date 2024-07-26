package uptime

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, struct {
		Error string `json:"error"`
	}{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("error write response data: %s", err)
	}
}
