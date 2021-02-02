package postservice

import (
	"context"
)

// PostStore is the interface for a post store
type PostStore interface {
	CreatePost(context.Context, *Post) error
	DeletePost(context.Context, string) error
	GetPostBySlug(context.Context, string) (*Post, error)
	UpdatePost(context.Context, *Post) error
	GetPostSummaries(context.Context, int) (*PostSummaryList, error)
	GetPostSummariesByTag(context.Context, int, string) (*PostSummaryList, error)
}
