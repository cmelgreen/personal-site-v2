package main

import (
	"context"
	"os"
	"time"

	"personal-site-v2/backend/server/aws"
	"personal-site-v2/backend/server/database"
	"personal-site-v2/backend/server/imageresizeservice"
	"personal-site-v2/backend/server/postservice"

	"github.com/go-chi/chi"
	// "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// PULL INTO YAML FILE
const (
	// Default timeout length
	timeout = 10

	// Default environment variable for serving and default port
	portEnvVar  = "PERSONAL_SITE_PORT"
	defaultPort = ":80"
	frontendDir = "/frontend/static"
	https       = true

	// Environment vars/files to check for AWS CLI & SSM configuration
	baseAWSRegion  = "AWS_REGION"
	baseAWSRoot    = "AWS_ROOT"
	baseConfigName = "base_config"
	baseConfigPath = "/go/app_data/"
	withEncrpytion = true

	// Path to serve api at
	apiRoot = "/api"
)

// Create router and environment then serve
func main() {
	// Setup Server
	ctx, cancelFn := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancelFn()

	s := newServer(ctx)

	// Setup DB for API
	dbConfig := database.DBConfigFromAWS{
		BaseAWSRegion:  baseAWSRegion,
		BaseAWSRoot:    baseAWSRoot,
		BaseConfigName: baseConfigName,
		BaseConfigPath: baseConfigPath,
		WithEncrpytion: withEncrpytion,
	}

	s.log.Println(dbConfig.ConfigString(ctx))

	// dbConfig := database.DBConfigFromValues{
	// 	Database: "postgres",
	// 	Host:  
	// 	Port:     "5432",
	// 	User:     "postgres",
	// 	Password: "postgres-personal-site",
	// }

	s.newDBConnection(ctx, dbConfig)

	// Add Backend API routes and utils
	richTextParser := &DraftJS{}

	postService := postservice.NewPostService(s.db, richTextParser)
	authApp := newFirebaseAuth("../credentials/firebase.json")

	s.mux.Use(cors.AllowAll().Handler)
	//s.mux.Use(middleware.Compress(5))

	s.mux.Get(apiRoot+"/post/{slug}", postService.GetPostHTTP())
	s.mux.Get(apiRoot+"/post-summaries", postService.GetPostSummariesHTTP())

	s.mux.Group(func(r chi.Router) {
		//r.Use(testMiddlware)
		r.Use(firebaseAuth(authApp))

		r.Post(apiRoot+"/post/", postService.CreatePostHTTP())
		r.Put(apiRoot+"/post/", postService.UpdatePostHTTP())
		r.Delete(apiRoot+"/post/{slug}", postService.DeletePostHTTP())
	})

	// Get port and serve
	port := os.Getenv(portEnvVar)
	if port == "" {
		port = defaultPort
	}

	s.mux.Post(apiRoot+"/img/", imageresizeservice.CreateImageHTTP("/", "test"))
	s.mux.Get(apiRoot+"/img/{img}", serveDynamicImage())
	s.mux.Get(apiRoot+"/status", status)
	s.log.Println("Serving:")
	s.log.Println("https: ", https)
	s.printRoutes()

	if https {
		cache := aws.NewSSMCache(true, "personal-site/certs")
		s.serveHTTPS(cache)
	} else {
		s.serve(port)
	}

}
