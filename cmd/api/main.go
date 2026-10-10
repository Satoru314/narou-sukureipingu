package main

import (
	"log"
	"net/http"
	"os"

	"narou-scraper/internal/db"
	"narou-scraper/internal/handler"
)

func main() {
	// Initialize database
	database, err := db.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Successfully connected to database")

	// Initialize handlers
	menuHandler := handler.NewMenuHandler(database)
	storyHandler := handler.NewStoryHandler(database)
	fetchAllHandler := handler.NewFetchAllHandler(database)

	// Setup routes
	http.HandleFunc("/api/menu", menuHandler.Handle)
	http.HandleFunc("/api/story", storyHandler.Handle)
	http.HandleFunc("/api/novel/fetch-all", fetchAllHandler.Handle)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
