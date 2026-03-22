package dto

import "time"

type User struct {
	FirstName     string `json:"first_name" binding:"required,safe_alphabets"`
	LastName      string `json:"last_name" binding:"required,safe_alphabets"`
	Rank          string `json:"rank" binding:"required,safe_alphabets"`
	SemanBook     string `json:"seman_book" binding:"required,safe_url"`
	Licence       string `json:"licence" binding:"required,safe_url"`
	LicenceType   string `json:"licence_type" binding:"required,safe_alphabets"`
	LicenceNumber string `json:"licence_number" binding:"required,safe_alphabets_with_numbers"`
	CompanyName   string `json:"company_name" binding:"required,no_html,max=1000"`
}

type FollowRequest struct {
	FolloweeID int `json:"followee_id" binding:"safe_int"`
}

type UnFollowRequest struct {
	FollowerID int `json:"follower_id" binding:"safe_int"`
}

type UpdateBioAndProfilePic struct {
	ProfilePicture string `json:"profile_picture" binding:"safe_url"`
	Bio            string `json:"bio"`
}

type UserProfileResponse struct {
	UserProfile UserProfile
	Cursor      int64
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
	Posts          []Post
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
