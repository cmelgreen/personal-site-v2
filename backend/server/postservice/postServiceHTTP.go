package postservice

import (
	"encoding/json"
	"net/http"
)

// CreatePostHTTP serves
func (p *PostService) CreatePostHTTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := p.createPost(r.Context(), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
}

// UpdatePostHTTP serves
func (p *PostService) UpdatePostHTTP() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		err := p.updatePost(r.Context(), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
}

// DeletePostHTTP serves
func (p *PostService) DeletePostHTTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := parseRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err = p.deletePost(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		w.WriteHeader(http.StatusOK)
	}
}

// GetPostHTTP serves
func (p *PostService) GetPostHTTP() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := parseRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		item, err := p.getPostBySlug(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		json.NewEncoder(w).Encode(item)
	}
}

// GetPostSummariesHTTP serves
func (p *PostService) GetPostSummariesHTTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := parseRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		item, err := p.getPostSummaries(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		json.NewEncoder(w).Encode(item)
	}
}
