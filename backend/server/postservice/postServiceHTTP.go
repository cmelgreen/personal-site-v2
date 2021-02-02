package postservice

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func parseRequest(r *http.Request, params httprouter.Params) (*PostRequest, error) {
	var request *PostRequest

	err := UnmarshalRequest(r, request)
	if err != nil {
		return nil, err
	}
	request.Slug = params.ByName("slug")

	return request, nil
}

func (p *PostService) getPostHTTP() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		request, err := parseRequest(r, params)
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

func (p *PostService) createPostHTTP() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		err := p.createPost(r.Context(), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (p *PostService) updatePostHTTP() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		err := p.updatePost(r.Context(), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (p *PostService) deletePostHTTP() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		request, err := parseRequest(r, params)
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

func (p *PostService) getPostSummariesHTTP() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		request, err := parseRequest(r, params)
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
