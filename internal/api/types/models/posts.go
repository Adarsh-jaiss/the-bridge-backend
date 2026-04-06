package models

type Post struct {
	ID          int64
	Type        string
	Content     string
	Attachments []string
	PollOptions []string
}
