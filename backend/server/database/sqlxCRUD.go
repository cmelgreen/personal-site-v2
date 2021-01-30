package database

import (
	"context"
	"fmt"

	"PersonalSite/backend/models"

	// Database driver for db interactions
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jmoiron/sqlx"
)


// CreatePost creates
func (db *Database) CreatePost(ctx context.Context, post *models.Post) error {
	queries := []string{
		`INSERT INTO post(title, slug, img, summary, content, raw_content) 
		VALUES (:title, :slug, :img, :summary. :content, :raw_content));`,

		`INSERT INTO tag(value)
		VALUES (UNNEST(:tags.value))
		ON CONFLICT DO NOTHING;`,

		`INSERT INTO post_TAG
		SELECT    p.id, t.id
		FROM      post p
		LEFT JOIN tag t
		ON        t.value = ANY(:tags.value)
		WHERE	  p.slug = :slug;`,
	}


	return nil
}

// DeletePost deletes
func (db *Database) DeletePost(ctx context.Context) error {
	return nil
}

// GetPostBySlug gets
func (db *Database) GetPostBySlug(ctx context.Context) error {
	return nil
}

// UpdatePostBySlug updates
func (db *Database) UpdatePostBySlug(ctx context.Context) error {
	return nil
}

// GetPostSummaries gets
func (db *Database) GetPostSummaries(ctx context.Context) error {
	return nil
}