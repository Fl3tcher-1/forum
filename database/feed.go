package database

import "time"

type PostFeed struct {
	ID       int
	Title    string
	Content  string
	Comments []string
	Likes    int
	Created  time.Time
	Category string
}
