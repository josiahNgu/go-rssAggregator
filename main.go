package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rssaggregator/internal/database"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Hello World")
	godotenv.Load()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT NOT DEFINED")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DBUrl NOT DEFINED")

	}
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error with db connection", err)
	}
	// now we have a db config that we can pass to our handler so that they have access to our DB
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	// we call it before listenAndServer() because that function will block
	go startScrapping(db, 10, time.Minute)
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	srv := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	v1Router := chi.NewRouter()
	v1Router.Get("/health", handlerReadiness)
	v1Router.Get("/error", handleErr)
	v1Router.Post("/user", apiCfg.handlerCreateUser)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))
	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feed", apiCfg.middlewareAuth(apiCfg.handlerGetFeeds))
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))
	router.Mount("/v1", v1Router)
	log.Printf("Server starting on Port %s", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port", port)
}
