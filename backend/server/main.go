package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"PersonalSite/backend/database"
	"PersonalSite/backend/utils"
)

// PULL INTO YAML FILE
const (
	// Default timeout length
	timeout = 10

	// Default environment variable for serving and default port
	portEnvVar  = "PORT"
	defaultPort = ":80"
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
	dbConfig := database.DBConfigFromAWS{
		BaseAWSRegion:  baseAWSRegion,
		BaseAWSRoot:    baseAWSRoot,
		BaseConfigName: baseConfigName,
		BaseConfigPath: baseConfigPath,
		WithEncrpytion: withEncrpytion,
	}

	s.newDBConnection(ctx, dbConfig)

	// Add Backend API routes and utils
	richTextEditor := &utils.DraftJS{}
	
	s.mux.GET(apiRoot+"/post/:slug", s.getPostByID())
	s.mux.POST(apiRoot+"/post/:slug", s.createPost(richTextEditor))
	s.mux.PUT(apiRoot+"/post/:slug", s.updatePost(richTextEditor))
	s.mux.DELETE(apiRoot+"/post/:slug", s.deletePost())
	s.mux.GET(apiRoot+"/post-summaries", s.getPostSummaries())

	// Get port and serve
	port := os.Getenv(portEnvVar)
	if port == "" {
		port = defaultPort
	}

	// TODO - SERVE HTTPS
	s.log.Fatal(http.ListenAndServe(port, s.mux))
}
