package database

import (
	"context"
	"fmt"

	"personal-site-v2/backend/server/postservice"
)

// CreatePost creates
func (db *Database) CreatePost(ctx context.Context, post *postservice.Post) error {
	queries := []string{
		`INSERT INTO post(title, slug, img, summary, category, content, raw_content) 
		VALUES (:title, :slug, :img, :summary, :category, :content, :raw_content);`,

		`INSERT INTO tag(value)
		VALUES (UNNEST(:tags ::::text[]))
		ON CONFLICT DO NOTHING;`,

		`INSERT INTO post_TAG
		SELECT    p.id, t.id
		FROM      post p
		LEFT JOIN tag t
		ON        t.value = ANY(:tags ::::text[])
		WHERE	  p.slug = :slug;`,
	}

	return db.TransactionCtx(ctx, queries, *post)
}

// DeletePost deletes
func (db *Database) DeletePost(ctx context.Context, slug string) error {
	query := "DELETE FROM post WHERE slug=$1;"
	_, err := db.ExecContext(ctx, query, slug)

	return err
}

// GetPostBySlug gets
func (db *Database) GetPostBySlug(ctx context.Context, slug string) (*postservice.Post, error) {
	var post postservice.Post

	query := "SELECT title, slug, img, summary, category, content, raw_content FROM post WHERE slug=$1;"
	err := db.GetContext(ctx, &post, query, slug)

	post.Tags = db.GetTagsBySlug(ctx, slug)

	return &post, err
}

// UpdatePost updates
func (db *Database) UpdatePost(ctx context.Context, post *postservice.Post) error {
	queries := []string{
		`UPDATE post SET (title, slug, img, summary, category, content, raw_content) =
			(:title, :slug, :img, :summary, :category, :content, :raw_content)
		WHERE slug = :slug;`,

		`DELETE FROM post_tag
		WHERE post_id in (
			SELECT id FROM post
			WHERE slug = :slug
		);`,

		`INSERT INTO tag(value)
		VALUES (UNNEST(:tags ::::text[]))
		ON CONFLICT DO NOTHING;`,

		`INSERT INTO post_tag
		SELECT    p.id, t.id
		FROM      post p
		LEFT JOIN tag t
		ON        t.value = ANY(:tags ::::text[])
		WHERE	  p.slug = :slug;`,
	}

	return db.TransactionCtx(ctx, queries, post)
}

// GetPostSummaries gets
func (db *Database) GetPostSummaries(ctx context.Context, limit int) (*postservice.PostSummaryList, error) {
	var posts []*postservice.PostSummary

	query := `SELECT title, slug, img, summary, category FROM post LIMIT $1;`
	err := db.SelectContext(ctx, &posts, query, limit)

	for _, post := range posts {
		post.Tags = db.GetTagsBySlug(ctx, post.Slug)
	}

	return &postservice.PostSummaryList{Posts: posts}, err
}

// GetPostSummariesByTag gets
func (db *Database) GetPostSummariesByTag(ctx context.Context, limit int, tag string) (*postservice.PostSummaryList, error) {
	var posts []*postservice.PostSummary

	query := `
	SELECT  title, slug, img, summary, category FROM post
	WHERE id in (
		SELECT post_id 
		  FROM post_to_tag
		  WHERE value = $2
	)
	LIMIT $1;`
	err := db.SelectContext(ctx, &posts, query, limit, tag)

	for _, post := range posts {
		post.Tags = db.GetTagsBySlug(ctx, post.Slug)
	}

	return &postservice.PostSummaryList{Posts: posts}, err
}

// GetTagsBySlug gets
func (db *Database) GetTagsBySlug(ctx context.Context, slug string) []string {
	var tags []string

	query := "SELECT value FROM post_to_tag WHERE slug=$1;"
	rows, err := db.QueryxContext(ctx, query, slug)
	if err != nil {
		fmt.Println(err)
		return tags
	}

	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			fmt.Println(err)
			return tags
		}
		tags = append(tags, s)
	}

	return tags
}

// TransactionCtx takes a context, slice of queries, and argument. Commits and returns nil if all queries
// execute succesffuly otherwise rollsback transaction and returns err
func (db *Database) TransactionCtx(ctx context.Context, queries []string, arg interface{}) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	for _, query := range queries {
		_, err = tx.NamedExecContext(ctx, query, arg)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
