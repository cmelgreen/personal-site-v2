package postservice

import (
	"io"
	"encoding/json"
	"context"
)

// PostService for posts
type PostService struct {
	store PostStore
	rtParser RichTextParser
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

type postUser = func(context.Context, *Post) error

func (p *PostService) getPostBySlug(ctx context.Context, r *PostRequest) (*Post, error) {
	post, err := p.store.GetPostBySlug(ctx, r.Slug)
	if err != nil {
		return nil, err
	}

	if r.Raw {
		post.Content = post.RawContent
		post.RawContent = ""
	}

	return post, nil
}

func (p *PostService) createPost(ctx context.Context, r io.Reader) error {
	return p.unmarshalPostAndUse(ctx, r, p.store.CreatePost)
}

func (p *PostService) updatePost(ctx context.Context, r io.Reader) error {
	return p.unmarshalPostAndUse(ctx, r, p.store.UpdatePost)
}

func (p *PostService) unmarshalPostAndUse(ctx context.Context, r io.Reader, f postUser) error {
	var post Post

	err := json.NewDecoder(r).Decode(&post)
	if err != nil {
		return err
	}

	html, err := p.rtParser.RichTextToHTML(post.Content)
	if err != nil {
		return nil
	}

	post.RawContent = post.Content
	post.Content = html

	err = f(ctx, &post)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) deletePost(ctx context.Context, r *PostRequest) error {
	err := p.store.DeletePost(ctx, r.Slug)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostService) getPostSummaries(ctx context.Context, r *PostRequest) (*PostSummaryList, error) {
	var postSummaries *PostSummaryList
	var err error

	if r.Num == 0 {
		r.Num = 10
	}

	if r.Tag != "" {
		postSummaries, err = p.store.GetPostSummariesByTag(ctx, r.Num, r.Tag)
	} else {
		postSummaries, err = p.store.GetPostSummaries(ctx, r.Num)
	}
	if err != nil {
		return nil, err
	}

	return postSummaries, nil
}

