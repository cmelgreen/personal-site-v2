package models

// Post is the main structure served and displayed
type Post struct {
	Title 		string 	`json:"title" db:"title"`
	Slug		string 	`json:"slug" db:"slug"`
	Img			string 	`json:"img" db:"img"`
	Summary		string 	`json:"summary" db:"summary"`
	Content 	string 	`json:"content" db:"content"`
	RawContent  string	`db:"raw_content"`
	Tags 		Tags	`json:"tags"`
}

// PostList is a list of Posts
type PostList struct {
	Posts 		[]*Post `json:"posts"`
}

// PostSummary is the summary information for a post
type PostSummary struct {
	Title		string	`json:"title" db:"title"`
	Slug		string 	`json:"slug" db:"slug"`
	Thumbnail	string	`json:"thumbnail" db:"thumbnail"`
	Summary		string	`json:"summary" db:"summary"`
	Tags		Tags	`json:"tags"`
}

// PostSummaryList is a list of PostSummaries
type PostSummaryList struct {
	Posts 		[]*PostSummary `json:"posts"`
}

// Tags are post tags 
type Tags struct {
	Values []string  `json:"value" db:"value"`
}