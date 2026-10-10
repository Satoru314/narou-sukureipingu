package db

import (
	"fmt"

	"narou-scraper/internal/models"
)

// SaveEpisodes inserts or updates multiple episodes
func (db *DB) SaveEpisodes(ncode string, episodes []models.EpisodeListItem) error {
	query := `
		INSERT INTO episodes (ncode, episode_number, subtitle, url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (ncode, episode_number) DO UPDATE SET
			subtitle = EXCLUDED.subtitle,
			url = EXCLUDED.url`

	for _, ep := range episodes {
		_, err := db.conn.Exec(query, ncode, ep.Number, ep.Subtitle, ep.URL)
		if err != nil {
			return fmt.Errorf("failed to save episode %d: %w", ep.Number, err)
		}
	}

	return nil
}

// UpdateEpisodeContent updates the content and posted_at of an episode
func (db *DB) UpdateEpisodeContent(ncode string, episodeNumber int, story *models.StoryResponse) error {
	query := `UPDATE episodes SET content = $1, posted_at = $2 WHERE ncode = $3 AND episode_number = $4`

	_, err := db.conn.Exec(query, story.Content, story.PostedAt, ncode, episodeNumber)
	if err != nil {
		return fmt.Errorf("failed to update episode content: %w", err)
	}

	return nil
}

// GetEpisodes retrieves all episodes for a novel
func (db *DB) GetEpisodes(ncode string) ([]*models.Episode, error) {
	query := `SELECT ncode, episode_number, subtitle, url, content, posted_at, created_at
	          FROM episodes WHERE ncode = $1 ORDER BY episode_number`

	rows, err := db.conn.Query(query, ncode)
	if err != nil {
		return nil, fmt.Errorf("failed to get episodes: %w", err)
	}
	defer rows.Close()

	var episodes []*models.Episode
	for rows.Next() {
		ep := &models.Episode{}
		err := rows.Scan(&ep.NCode, &ep.EpisodeNumber, &ep.Subtitle, &ep.URL, &ep.Content, &ep.PostedAt, &ep.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan episode: %w", err)
		}
		episodes = append(episodes, ep)
	}

	return episodes, nil
}
