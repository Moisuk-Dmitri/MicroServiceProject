package blogmodel

import "time"

type Post struct {
	ID        string
	AuthorID  string
	Title     string
	Content   string
	CreatedAt time.Time
}
