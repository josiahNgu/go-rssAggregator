package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssaggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type paramteres struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramteres{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing Json %v", err))
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       params.URL,
		UserUuid:  user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user %v", err))
	}
	respondWithJson(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	userFeeds, err := apiCfg.DB.GetUserFeeds(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error getting feed %v", err))
	}
	fmt.Println(userFeeds)
	respondWithJson(w, 200, databaseFeedsToFeeds(userFeeds))
}
