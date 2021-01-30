package database

import (
	"context"

	"PersonalSite/backend/models"
)

// CreatePost creates
func (db *Database) CreatePost(ctx context.Context, post *models.Post) error {
	queries := []string{
		`INSERT INTO post(title, slug, img, summary, content, raw_content) 
		VALUES (:title, :slug, :img, :summary. :content, :raw_content);`,

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

	return db.TransactionCtx(ctx, queries, post)
}

// DeletePost deletes
func (db *Database) DeletePost(ctx context.Context, slug string) error {
	query := "DELETE FROM post WHERE slug=$1;"
	_, err := db.ExecContext(ctx, query, slug)

	return err
}

// GetPostBySlug gets
func (db *Database) GetPostBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var post models.Post

	query := "SELECT id, title, summary, content FROM post WHERE slug=$1;"
	err := db.GetContext(ctx, &post, query, slug)

	post.Tags = GetTagsBySlug(ctx, slug)

	return &post, err
}

// GetPostRawBySlug gets
func (db *Database) GetPostRawBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var post models.Post

	query := "SELECT id, title, summary, content_raw as content FROM post WHERE slug=$1;"
	err := db.GetContext(ctx, &post, query, slug)

	var tags models.Tags
	query = "SELECT value FROM post_to_tag WHERE slug=$1;"
	err = db.GetContext(ctx, &tags, query, slug)

	post.Tags = tags

	return &post, err
}

// UpdatePostBySlug updates
func (db *Database) UpdatePost(ctx context.Context, post *models.Post) error {
	queries := []string{
		`UPDATE post SET (title, slug, img, summary, content, raw_content) =
			(:title, :slug, :img, :summary. :content, :raw_content)
		WHERE slug = :slug;`,

		`DELETE FROM post_tag
		WHERE post_id in (
			SELECT id FROM post
			WHERE slug = :slug
		);`,

		`INSERT INTO tag(value)
		VALUES (UNNEST(:tags.value))
		ON CONFLICT DO NOTHING;`,

		`INSERT INTO post_tag
		SELECT    p.id, t.id
		FROM      post p
		LEFT JOIN tag t
		ON        t.value = ANY(:tags.value)
		WHERE	  p.slug = 'aaa';`,
	}

	return db.TransactionCtx(ctx, queries, post)
}

// GetPostSummaries gets
func (db *Database) GetPostSummaries(ctx context.Context, limit int) (*models.PostSummaryList, error) {
	var posts []*models.PostSummary

	query := `SELECT id, title, summary FROM post LIMIT $1;`
	err := db.SelectContext(ctx, &posts, query, limit)

	for _, post := range posts {
		post.Tags = *db.GetTagsBySlug(ctx, post.Slug)
	}

	return &models.PostSummaryList{Posts: posts}, err
}

// GetPostSummariesByTag gets
func (db *Database) GetPostSummariesByTag(ctx context.Context, limit int, tag string) (*models.PostSummaryList, error) {
	var posts []*models.PostSummary

	query := `
	SELECT * FROM post
	WHERE id in (
		SELECT post_id 
		  FROM post_to_tag
		  WHERE value = $2
	)
	LIMIT $1;`
	err := db.SelectContext(ctx, &posts, query, limit, tag)

	return &models.PostSummaryList{Posts: posts}, err
}

// GetTagsBySlug gets
func (db *Database) GetTagsBySlug(ctx context.Context, slug string) *models.Tags {
	var tags models.Tags

	query := "SELECT value FROM post_to_tag WHERE slug=$1;"
	db.GetContext(ctx, &tags, query, slug)

	return &tags
}

// TransactionCtx takes a context, slice of queries, and argument. Commits and returns nil if all queries
// execute succesffuly otherwise rollsback transaction and returns err
func (db *Database) TransactionCtx(ctx context.Context, queries []string, arg interface{}) error {
	txx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, query := range queries {
		_, err = txx.NamedExecContext(ctx, query, arg)
		if err != nil {
			txx.Rollback()
			return err
		}
	}

	txx.Commit()
	return nil
}