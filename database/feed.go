package database

import "encoding/json"

type PostFeed struct {
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	Created  string `json:"created,omitempty"`
	Category string `json:"category,omitempty"`
	// Comments []string
}

func (p PostFeed) MarshallJSON() ([]byte, error) {
	return json.Marshal(p)
}