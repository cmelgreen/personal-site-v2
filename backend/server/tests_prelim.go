package main

import (
	"context"

	"PersonalSite/backend/models"
)

func createDummyPost(ctx context.Context, s *Server) {
	testPosts := []models.Post{
		{
			Title: "Automating Dev Environment Setup",
			Slug: "dev-env-setup",
			Summary: "Next time you blow up your environment restore it in one command with Ansible",
			Category: "DevOps",
			Img: "ansible.png",
		},{
			Title: "Reproducible CI/CD Pipelines",
			Slug: "reproducible-ci-cd",
			Summary: "Use Pipelines and Jenkins-configuration-as-Code to skip the manual setup",
			Category: "DevOps",
			Img: "jenkins.png",
		},{
			Title: "React & Go I: A simple server",
			Slug: "react-golang-simple",
			Summary: "Getting your feet wet using Go and React together",
			Category: "Backend",
			Img: "go-with-react.png",
		},{
			Title: "React & Go II: React Router and APIs",
			Slug: "react-golang-advanced",
			Summary: "Building complex interactions with React and Go",
			Category: "Backend",
			Img: "go-with-react.png",
		},
	}

	testTags := []models.Tags{
		{Values: []string{"ansible", "linux"}},
		{Values: []string{"jenkins", "ci/cd"}},
		{Values: []string{"go", "react"}},
		{Values: []string{"go", "react"}},
	}

	for i := range testPosts {
		testPosts[i].Tags = testTags[i]
		err := s.db.CreatePost(ctx, &testPosts[i])
		if err != nil {
			s.log.Println(err)
		}
	}
}



    