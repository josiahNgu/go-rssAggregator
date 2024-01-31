package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to Marshall Json response: %v", payload)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error", msg)
	}
	// provide a syntax on how to unravel the json
	// {error: "something here"}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}
