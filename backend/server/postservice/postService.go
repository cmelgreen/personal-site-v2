package postservice

import (
	"io"
	"encoding/json"
	"context"
	"log"
)



// NewPostService is the entry point
func NewPostService(store PostStore, rtParser RichTextParser) *PostService {
	return &PostService{
		store: store,
		rtParser: rtParser,
	}
}

func (p *PostService) getPostBySlug(ctx context.Context, r *PostRequest) (*Post, error) {
	post, err := p.store.GetPostBySlug(ctx, r.Slug)
	if err != nil {
		return nil, err
	}

	log.Println("Post: ", post)

	if r.Raw {
		log.Println("Raw: ", post.RawContent)
		log.Println("Con: ", post.Content)
		post.Content = post.RawContent
		post.RawContent = ""
	}

	log.Println(post)

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
		return err
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

