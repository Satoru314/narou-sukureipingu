package scraper

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"narou-scraper/internal/models"
)

type StoryScraper struct {
	client *http.Client
}

func NewStoryScraper() *StoryScraper {
	return &StoryScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ScrapeStory scrapes a single episode
func (s *StoryScraper) ScrapeStory(episodeURL string) (*models.StoryResponse, error) {
	// Build full URL if needed
	fullURL := episodeURL
	if !strings.HasPrefix(episodeURL, "http") {
		fullURL = "https://ncode.syosetu.com" + episodeURL
	}

	doc, err := s.fetchPage(fullURL)
	if err != nil {
		return nil, err
	}

	subtitle := strings.TrimSpace(doc.Find(".p-novel__title").Text())

	var contentParts []string
	doc.Find(".p-novel__text p").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		if text != "" {
			contentParts = append(contentParts, text)
		}
	})

	content := strings.Join(contentParts, "\n")

	episodeNum := s.extractEpisodeNumber(episodeURL)

	// Try to extract posted date
	var postedAt *time.Time
	postedDateStr := doc.Find("#novel_no").Text()
	if postedDateStr != "" {
		if parsed := s.parseDate(postedDateStr); !parsed.IsZero() {
			postedAt = &parsed
		}
	}

	return &models.StoryResponse{
		EpisodeNumber: episodeNum,
		Subtitle:      subtitle,
		Content:       content,
		PostedAt:      postedAt,
	}, nil
}

func (s *StoryScraper) fetchPage(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *StoryScraper) extractEpisodeNumber(url string) int {
	// Extract episode number from URL like "/n9619hu/123/" or "https://ncode.syosetu.com/n9619hu/123/"
	re := regexp.MustCompile(`/(\d+)/?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 1
	}

	num := 1
	fmt.Sscanf(matches[1], "%d", &num)
	return num
}

func (s *StoryScraper) parseDate(dateStr string) time.Time {
	// Parse date from string like "2023/01/15 12:34"
	parts := strings.Fields(dateStr)
	for _, part := range parts {
		if strings.Contains(part, "/") {
			dateParts := strings.Split(part, "/")
			if len(dateParts) == 3 {
				year, _ := strconv.Atoi(dateParts[0])
				month, _ := strconv.Atoi(dateParts[1])
				day, _ := strconv.Atoi(dateParts[2])
				if year > 0 && month > 0 && day > 0 {
					return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
				}
			}
		}
	}
	return time.Time{}
}
