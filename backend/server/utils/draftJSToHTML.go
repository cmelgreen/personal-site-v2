package utils

import (
	"encoding/json"

	"github.com/ejilay/draftjs"
)

// DraftJS hangs richTextToHTML method for Draft-JS output 
type DraftJS struct{}

// RichTextToHTML converts Draft-JS state string to HTML
func (d *DraftJS) RichTextToHTML(rawState string) (string, error) {
	contentState := draftjs.ContentState{}
	err := json.Unmarshal([]byte(rawState), &contentState)
	if err != nil {
		return "", err
	}

	config := draftjs.NewDefaultConfig()

	html := draftjs.Render(&contentState, config)

	return html, nil

}