package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World")
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT NOT DEFINED")
	}
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
	router.Mount("/v1", v1Router)
	log.Printf("Server starting on Port %s", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port", port)
}
