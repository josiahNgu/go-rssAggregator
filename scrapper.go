package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"rssaggregator/internal/database"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// long running job that run in the background as server running
func startScrapping(
	db *database.Queries,
	concurrency int,
	timeout time.Duration,
) {
	fmt.Printf("scrapping on %v go routines every %s duration", concurrency, timeout)
	ticker := time.NewTicker(timeout)
	// if timeout is 1 min, we are saying to run this for loop every 1 min
	// context.background is the global context, its what you used if you have no access to local context
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeeds:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	// this decrement the counter by 1 so we can call scrapeFeed x times and line 33 will execute when everything is done
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking as fetched:", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubDate,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("failed to create post %v to DB %v", feed.ID, err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
