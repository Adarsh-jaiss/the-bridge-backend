package models

import "time"

type User struct {
	ID            int64
	FirstName     string
	LastName      string
	Rank          string
	IsVerified    bool
	SemanBook     string
	CompanyName   string
	Licence       string
	LicenceType   string
	LicenceNumber string
}

type UserProfile struct {
	FirstName      string
	LastName       string
	ProfilePicture *string
	Rank           string
	Bio            *string
	FollowersCount int64
	FollowingCount int64
	TotalPostCount int64
}

type Post struct {
	ID          int64
	Content     *string   // text content of the post
	Attachments []string // list of images/video url's
	Likes       int64
	Comments    int64
	Reposts     int64
	Bookmarks   int64
	CreatedAt   time.Time
}
