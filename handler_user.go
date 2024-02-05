package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssaggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type paramteres struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramteres{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing Json %v", err))
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user %v", err))
	}
	respondWithJson(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJson(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{UserID: user.ID, Limit: 10})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Coundn't find user posts %v", user.ID))
	}
	respondWithJson(w, 200, databasePostsToPosts(posts))
}
