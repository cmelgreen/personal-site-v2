package postservice

import (
	"context"
)

// PostService for posts
type PostService struct {
	store PostStore
	rtParser RichTextParser
}

// PostStore is the interface for a post store
type PostStore interface {
	CreatePost(context.Context, *Post) error
	DeletePost(context.Context, string) error
	GetPostBySlug(context.Context, string) (*Post, error)
	UpdatePost(context.Context, *Post) error
	GetPostSummaries(context.Context, int) (*PostSummaryList, error)
	GetPostSummariesByTag(context.Context, int, string) (*PostSummaryList, error)
}

// PostRequest represents the request for a post
type PostRequest struct {
	Slug   string 
	Num    int    `request:"numPosts"`
	Raw    bool   `request:"raw"`
	SortBy string `request:"sortBy"`
	Tag    string `request:"tag"`
}

// RichTextParser is interface for converting Rich Text Editor output to HTML
type RichTextParser interface {
	RichTextToHTML(string) (string, error)
}

type postUser func(context.Context, *Post) error

// Post is the main structure served and displayed
type Post struct {
	Title 		string 	`json:"title" db:"title"`
	Slug		string 	`json:"slug" db:"slug"`
	Img			string 	`json:"img" db:"img"`
	Summary		string 	`json:"summary" db:"summary"`
	Category	string 	`json:"category" db:"category"`
	Content 	string 	`json:"content" db:"content"`
	RawContent  string	`db:"raw_content"`
	Tags 		[]string	`json:"tags"`
}

// PostList is a list of Posts
type PostList struct {
	Posts 		[]*Post `json:"posts"`
}

// PostSummary is the summary information for a post
type PostSummary struct {
	Title		string	`json:"title" db:"title"`
	Slug		string 	`json:"slug" db:"slug"`
	Img			string	`json:"img" db:"img"`
	Summary		string	`json:"summary" db:"summary"`
	Category	string 	`json:"category" db:"category"`
	Tags		[]string `json:"tags"`
}

// PostSummaryList is a list of PostSummaries
type PostSummaryList struct {
	Posts 		[]*PostSummary `json:"posts"`
}