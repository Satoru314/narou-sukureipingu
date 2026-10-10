package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"narou-scraper/internal/db"
	"narou-scraper/internal/models"
	"narou-scraper/internal/scraper"
)

type FetchAllHandler struct {
	db           *db.DB
	menuScraper  *scraper.MenuScraper
	storyScraper *scraper.StoryScraper
}

func NewFetchAllHandler(database *db.DB) *FetchAllHandler {
	return &FetchAllHandler{
		db:           database,
		menuScraper:  scraper.NewMenuScraper(),
		storyScraper: scraper.NewStoryScraper(),
	}
}

type FetchAllResponse struct {
	NCode         string `json:"ncode"`
	Title         string `json:"title"`
	TotalEpisodes int    `json:"total_episodes"`
	FetchedCount  int    `json:"fetched_count"`
	Message       string `json:"message"`
}

func (h *FetchAllHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get ncode from query parameter
	ncode := r.URL.Query().Get("ncode")
	if ncode == "" {
		http.Error(w, "ncode query parameter is required", http.StatusBadRequest)
		return
	}

	// Build URL from ncode
	novelURL := fmt.Sprintf("https://ncode.syosetu.com/%s", ncode)

	// Step 1: Fetch and save menu (table of contents)
	log.Printf("Fetching menu for %s...", ncode)
	menuData, err := h.menuScraper.ScrapeMenu(novelURL)
	if err != nil {
		http.Error(w, "Failed to scrape menu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.db.SaveNovel(menuData, ncode); err != nil {
		http.Error(w, "Failed to save novel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.db.SaveEpisodes(ncode, menuData.Episodes); err != nil {
		http.Error(w, "Failed to save episodes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Menu saved. Total episodes: %d", len(menuData.Episodes))

	// Get already fetched episodes from database
	existingEpisodes, err := h.db.GetEpisodes(ncode)
	if err != nil {
		log.Printf("Warning: Failed to get existing episodes: %v", err)
		existingEpisodes = []*models.Episode{}
	}

	// Create a map of already fetched episodes
	fetchedMap := make(map[int]bool)
	for _, ep := range existingEpisodes {
		if ep.Content != nil && *ep.Content != "" {
			fetchedMap[ep.EpisodeNumber] = true
		}
	}

	// Step 2: Fetch all episodes with 5-second delay (skip already fetched)
	fetchedCount := 0
	skippedCount := 0
	for i, episode := range menuData.Episodes {
		// Skip if already fetched
		if fetchedMap[episode.Number] {
			log.Printf("Episode %d/%d already fetched, skipping: %s", i+1, len(menuData.Episodes), episode.Subtitle)
			skippedCount++
			continue
		}

		episodeURL := fmt.Sprintf("https://ncode.syosetu.com%s", episode.URL)

		log.Printf("Fetching episode %d/%d: %s", i+1, len(menuData.Episodes), episode.Subtitle)

		storyData, err := h.storyScraper.ScrapeStory(episodeURL)
		if err != nil {
			log.Printf("Failed to scrape episode %d: %v", episode.Number, err)
			continue
		}

		if err := h.db.UpdateEpisodeContent(ncode, episode.Number, storyData); err != nil {
			log.Printf("Failed to save episode %d: %v", episode.Number, err)
			continue
		}

		fetchedCount++
		log.Printf("Episode %d/%d saved successfully", i+1, len(menuData.Episodes))

		// Wait 5 seconds before next request (except for the last one)
		if i < len(menuData.Episodes)-1 {
			time.Sleep(5 * time.Second)
		}
	}

	// Return response
	response := FetchAllResponse{
		NCode:         ncode,
		Title:         menuData.Title,
		TotalEpisodes: len(menuData.Episodes),
		FetchedCount:  fetchedCount,
		Message:       fmt.Sprintf("Successfully fetched %d new episodes, skipped %d already fetched (total: %d)", fetchedCount, skippedCount, len(menuData.Episodes)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
