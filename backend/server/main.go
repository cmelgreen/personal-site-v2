package main

import (
	"context"
	"os"
	"time"

	"PersonalSite/backend/database"

	
)

// PULL INTO YAML FILE
const (
	// Default timeout length
	timeout = 10

	// Default environment variable for serving and default port
	portEnvVar  = "PERSONAL_SITE_PORT"
	defaultPort = ":8080"
	frontendDir = "/frontend/static"

	// Environment vars/files to check for AWS CLI & SSM configuration
	baseAWSRegion  = "AWS_REGION"
	baseAWSRoot    = "AWS_ROOT"
	baseConfigName = "base_config"
	baseConfigPath = "./app_data/"
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
	// dbConfig := database.DBConfigFromAWS{
	// 	BaseAWSRegion:  baseAWSRegion,
	// 	BaseAWSRoot:    baseAWSRoot,
	// 	BaseConfigName: baseConfigName,
	// 	BaseConfigPath: baseConfigPath,
	// 	WithEncrpytion: withEncrpytion,
	// }

	dbConfig := database.DBConfigFromValues{
		Database: "postgres",
		Host: "localhost",
		Port: "5432",
		User: "postgres",
		Password: "postgres",
	}

	s.newDBConnection(ctx, dbConfig)

	// Add Backend API routes and utils
	richTextEditor := &DraftJS{}
	
	s.mux.GET(apiRoot+"/post/:slug", s.getPostBySlug())
	s.mux.POST(apiRoot+"/post", s.createPost(richTextEditor))
	s.mux.PUT(apiRoot+"/post", s.updatePost(richTextEditor))
	s.mux.DELETE(apiRoot+"/post/:slug", s.deletePost())
	s.mux.GET(apiRoot+"/post-summaries", s.getPostSummaries())

	// Get port and serve
	port := os.Getenv(portEnvVar)
	if port == "" {
		port = defaultPort
	}

	s.mux.GET(apiRoot+"/img/:img", serveDynamicImage())
	s.log.Println("Serving")

	//createDummyPost(ctx, s)

	s.serveCORSEnabled(port)
}
