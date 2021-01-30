package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"PersonalSite/backend/models"
	"PersonalSite/backend/utils"

	"github.com/julienschmidt/httprouter"
)

//go:generate go run utils/requestBinding/bindingGenerator.go -f postController.go -out APIBindings.go

// PostRequest represents the request for a post
type PostRequest struct {
	Num 		int		`request:"numPosts"`
	Raw			bool	`request:"raw"`
	SortBy		string	`request:"sortBy"`
	FilterBy	string	`request:"filterBy"`
}

// RichTextHandler is interface for converting Rich Text Editor output to HTML
type RichTextHandler interface{
	RichTextToHTML(string) (string, error)
}

func writeJSON(w http.ResponseWriter, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func unwrapBool(s string) bool {
	value, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return value
}

func (s *Server) getPostByID() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var request PostRequest 
		err := utils.UnmarshalRequest(r, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var post models.Post 


		if request.Raw {
			post, err = s.db.QueryPostRaw(r.Context(), request.Title)
		} else {
			post, err = s.db.QueryPost(r.Context(), request.Title)
		}

		if err != nil {
			s.log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
			// IMPLEMENT ERROR HANDLING
		}

		writeJSON(w, post)
	}
}

func (s *Server) createPost(richText RichTextHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var post models.Post

		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			s.log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
			// ADD ERROR HANDLING
		}

		html, err := richText.RichTextToHTML(post.Content)
		if err != nil {
			s.log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		post.RawContent = post.Content
		post.Content = html

		err = s.db.InsertPost(r.Context(), post)
		if err != nil {
			s.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) updatePost(richText RichTextHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var post models.Post

		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			s.log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
			// ADD ERROR HANDLING
		}

		html, err := richText.RichTextToHTML(post.Content)
		if err != nil {
			s.log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		post.RawContent = post.Content
		post.Content = html

		err = s.db.UpdatePost(r.Context(), post)
		if err != nil {
			s.log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) deletePost()  httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var request PostRequest 
		err := utils.UnmarshalRequest(r, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = s.db.DeletePost(r.Context(), request.Title)
		if err != nil {
			s.log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
			// IMPLEMENT ERROR HANDLING
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) getPostSummaries() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var request PostRequest
		err := utils.UnmarshalRequest(r, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		postSummaries, err := s.db.QueryPostSummaries(r.Context(), request.Num)
		if err != nil {
			s.log.Println(err)
		}

		writeJSON(w, postSummaries)
	}
}
