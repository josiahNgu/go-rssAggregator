package main

import (
	"fmt"
	"net/http"
	"rssaggregator/internal/auth"
	"rssaggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

// middlewareAuth is a function on our apiconfig, take authHandler as input and return Handler Func as output
func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Error with authorization %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("couldn't get user %v", err))
			return
		}
		handler(w, r, user)
	}
}
