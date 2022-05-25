package database

type PostFeed struct {
	ID      int
	Title   string
	Content string
	// Comments []string
	Likes    int
	Dislikes int
	Created  string
	Category string
}
