package db

import (
	"fmt"

	"narou-scraper/internal/models"
)

// SaveNovel inserts or updates a novel
func (db *DB) SaveNovel(novel *models.MenuResponse, ncode string) error {
	query := `
		INSERT INTO novels (ncode, title, author, author_url, description, novel_url, episode_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (ncode) DO UPDATE SET
			title = EXCLUDED.title,
			author = EXCLUDED.author,
			author_url = EXCLUDED.author_url,
			description = EXCLUDED.description,
			novel_url = EXCLUDED.novel_url,
			episode_count = EXCLUDED.episode_count,
			updated_at = CURRENT_TIMESTAMP`

	_, err := db.conn.Exec(query, ncode, novel.Title, novel.Author, novel.AuthorURL, novel.Description, novel.NovelURL, novel.EpisodeCount)
	if err != nil {
		return fmt.Errorf("failed to save novel: %w", err)
	}

	return nil
}

// GetNovelByNCode retrieves a novel by its ncode
func (db *DB) GetNovelByNCode(ncode string) (*models.Novel, error) {
	query := `SELECT ncode, title, author, author_url, description, novel_url, episode_count, created_at, updated_at
	          FROM novels WHERE ncode = $1`

	novel := &models.Novel{}
	err := db.conn.QueryRow(query, ncode).Scan(
		&novel.NCode, &novel.Title, &novel.Author, &novel.AuthorURL,
		&novel.Description, &novel.NovelURL, &novel.EpisodeCount,
		&novel.CreatedAt, &novel.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get novel: %w", err)
	}

	return novel, nil
}
