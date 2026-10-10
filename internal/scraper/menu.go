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

type MenuScraper struct {
	client *http.Client
}

func NewMenuScraper() *MenuScraper {
	return &MenuScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ScrapeMenu scrapes the table of contents for a novel
func (s *MenuScraper) ScrapeMenu(novelURL string) (*models.MenuResponse, error) {
	// Extract ncode from URL
	ncode, err := extractNCode(novelURL)
	if err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("https://ncode.syosetu.com/%s", ncode)

	// Fetch first page
	doc, err := s.fetchPage(baseURL)
	if err != nil {
		return nil, err
	}

	// Get basic info
	title := strings.TrimSpace(doc.Find(".p-novel__title").Text())
	author := strings.TrimSpace(doc.Find(".p-novel__author a").Text())

	// Get author URL
	authorURL := ""
	if href, exists := doc.Find(".p-novel__author a").Attr("href"); exists {
		authorURL = href
	}

	// Get description
	description := strings.TrimSpace(doc.Find("#novel_ex").Text())

	// Get total pages
	totalPages := s.getTotalPages(doc)

	var allEpisodes []models.EpisodeListItem
	currentEpisodeNum := 1

	// Scrape all pages
	for page := 1; page <= totalPages; page++ {
		var pageDoc *goquery.Document
		if page == 1 {
			pageDoc = doc
		} else {
			pageURL := fmt.Sprintf("%s?p=%d", baseURL, page)
			pageDoc, err = s.fetchPage(pageURL)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
			}
			time.Sleep(5 * time.Second) // Be respectful to the server
		}

		episodes := s.extractEpisodes(pageDoc, ncode, currentEpisodeNum)
		allEpisodes = append(allEpisodes, episodes...)
		currentEpisodeNum += len(episodes)
	}

	return &models.MenuResponse{
		NCode:        ncode,
		Title:        title,
		Author:       author,
		AuthorURL:    authorURL,
		Description:  description,
		NovelURL:     baseURL,
		EpisodeCount: len(allEpisodes),
		Episodes:     allEpisodes,
	}, nil
}

func (s *MenuScraper) fetchPage(url string) (*goquery.Document, error) {
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

func (s *MenuScraper) getTotalPages(doc *goquery.Document) int {
	lastPageHref, exists := doc.Find(".c-pager__item--last").Attr("href")
	if !exists {
		return 1 // Single page
	}

	// Extract page number from href like "/n9619hu/?p=5"
	re := regexp.MustCompile(`\?p=(\d+)`)
	matches := re.FindStringSubmatch(lastPageHref)
	if len(matches) < 2 {
		return 1
	}

	pages, err := strconv.Atoi(matches[1])
	if err != nil {
		return 1
	}

	return pages
}

func (s *MenuScraper) extractEpisodes(doc *goquery.Document, ncode string, startNum int) []models.EpisodeListItem {
	var episodes []models.EpisodeListItem
	episodeNum := startNum

	doc.Find(".p-eplist__subtitle").Each(func(i int, sel *goquery.Selection) {
		subtitle := strings.TrimSpace(sel.Text())
		episodeURL := fmt.Sprintf("/%s/%d/", ncode, episodeNum)

		episodes = append(episodes, models.EpisodeListItem{
			Number:   episodeNum,
			Subtitle: subtitle,
			URL:      episodeURL,
		})
		episodeNum++
	})

	return episodes
}

// ExtractNCodeFromURL extracts ncode from URL (exported for use in handlers)
func ExtractNCodeFromURL(url string) (string, error) {
	return extractNCode(url)
}

func extractNCode(url string) (string, error) {
	// Extract ncode from URL like "https://ncode.syosetu.com/n9619hu" or "n9619hu"
	re := regexp.MustCompile(`([nN]\d+[a-zA-Z]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid novel URL: %s", url)
	}
	return strings.ToLower(matches[1]), nil
}
