package models

import "time"

// Novel represents a novel from "Shosetsuka ni Narou"
type Novel struct {
	NCode        string    `json:"ncode" db:"ncode"`
	Title        string    `json:"title" db:"title"`
	Author       string    `json:"author" db:"author"`
	AuthorURL    *string   `json:"author_url,omitempty" db:"author_url"`
	Description  *string   `json:"description,omitempty" db:"description"`
	NovelURL     string    `json:"novel_url" db:"novel_url"`
	EpisodeCount int       `json:"episode_count" db:"episode_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Episode represents a single episode/chapter of a novel
type Episode struct {
	NCode         string     `json:"ncode" db:"ncode"`
	EpisodeNumber int        `json:"episode_number" db:"episode_number"`
	Subtitle      string     `json:"subtitle" db:"subtitle"`
	URL           string     `json:"url" db:"url"`
	Content       *string    `json:"content,omitempty" db:"content"`
	PostedAt      *time.Time `json:"posted_at,omitempty" db:"posted_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// MenuResponse is the response for /api/menu
type MenuResponse struct {
	NCode        string            `json:"ncode"`
	Title        string            `json:"title"`
	Author       string            `json:"author"`
	AuthorURL    string            `json:"author_url,omitempty"`
	Description  string            `json:"description,omitempty"`
	NovelURL     string            `json:"novel_url"`
	EpisodeCount int               `json:"episode_count"`
	Episodes     []EpisodeListItem `json:"episodes"`
}

// EpisodeListItem represents an episode in the menu list
type EpisodeListItem struct {
	Number   int    `json:"number"`
	Subtitle string `json:"subtitle"`
	URL      string `json:"url"`
}

// StoryResponse is the response for /api/story
type StoryResponse struct {
	EpisodeNumber int        `json:"episode_number"`
	Subtitle      string     `json:"subtitle"`
	Content       string     `json:"content"`
	PostedAt      *time.Time `json:"posted_at,omitempty"`
}
