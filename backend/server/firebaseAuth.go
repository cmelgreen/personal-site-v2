package main

import (
	"fmt"
	"context"
	"net/http"
	"strings"
  
	firebase "firebase.google.com/go"
	//"firebase.google.com/go/auth"
  
	"google.golang.org/api/option"
)

func newFirebaseAuth(credFile string) *firebase.App {
	opt := option.WithCredentialsFile(credFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil
	}

	return app
}

func testMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Println("in middlware")

		next.ServeHTTP(w, r)
	})
}

func firebaseAuth(app *firebase.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			client, err := app.Auth(context.Background())

			fmt.Println("REQUEST FOR RESTRICTED CONTENT")

			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			idToken := parseAuthToken(r)
			_, err = client.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
        		return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseAuthToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

