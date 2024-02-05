package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssaggregator/internal/database"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type paramteres struct {
		FeedID uuid.UUID `json:"feedId"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramteres{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing Json %v", err))
	}
	follow, err := apiCfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user %v", err))
	}
	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(follow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	userFollows, err := apiCfg.DB.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error getting feed follows %v", err))
	}
	respondWithJson(w, 200, databaseFeedFollowToFeedFollows(userFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIdStr := chi.URLParam(r, "feedFollowID")
	feedId, err := uuid.Parse(feedIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Failed to parse to uuid %v", err))
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{ID: feedId, UserID: user.ID})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("couldn't delete feed follow %v", err))
	}
	respondWithJson(w, 200, struct{}{})
}
