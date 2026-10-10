package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"narou-scraper/internal/db"
	"narou-scraper/internal/scraper"
)

type MenuHandler struct {
	db      *db.DB
	scraper *scraper.MenuScraper
}

func NewMenuHandler(database *db.DB) *MenuHandler {
	return &MenuHandler{
		db:      database,
		scraper: scraper.NewMenuScraper(),
	}
}

func (h *MenuHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// Build URL
	novelURL := fmt.Sprintf("https://ncode.syosetu.com/%s", ncode)

	// Scrape menu
	menuData, err := h.scraper.ScrapeMenu(novelURL)
	if err != nil {
		http.Error(w, "Failed to scrape menu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save to database
	if err := h.db.SaveNovel(menuData, ncode); err != nil {
		http.Error(w, "Failed to save novel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.db.SaveEpisodes(ncode, menuData.Episodes); err != nil {
		http.Error(w, "Failed to save episodes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menuData)
}
