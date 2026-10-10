package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"narou-scraper/internal/db"
	"narou-scraper/internal/scraper"
)

type StoryHandler struct {
	db      *db.DB
	scraper *scraper.StoryScraper
}

func NewStoryHandler(database *db.DB) *StoryHandler {
	return &StoryHandler{
		db:      database,
		scraper: scraper.NewStoryScraper(),
	}
}

func (h *StoryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get ncode and episode_number from query parameters
	ncode := r.URL.Query().Get("ncode")
	episodeNumStr := r.URL.Query().Get("episode_number")

	if ncode == "" || episodeNumStr == "" {
		http.Error(w, "ncode and episode_number query parameters are required", http.StatusBadRequest)
		return
	}

	episodeNum, err := strconv.Atoi(episodeNumStr)
	if err != nil {
		http.Error(w, "episode_number must be a valid integer", http.StatusBadRequest)
		return
	}

	// Build URL
	episodeURL := fmt.Sprintf("https://ncode.syosetu.com/%s/%d/", ncode, episodeNum)

	// Scrape story
	storyData, err := h.scraper.ScrapeStory(episodeURL)
	if err != nil {
		http.Error(w, "Failed to scrape story: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update episode content in database
	if err := h.db.UpdateEpisodeContent(ncode, episodeNum, storyData); err != nil {
		http.Error(w, "Failed to save story: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storyData)
}
