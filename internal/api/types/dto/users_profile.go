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
	UserProfile UserProfile `json:"user_profile"`
	Cursor      int64       `json:"cursor"`
}

type UserProfile struct {
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	ProfilePicture *string `json:"profile_picture"`
	Rank           string  `json:"rank"`
	Bio            *string `json:"bio"`
	FollowersCount int64   `json:"followers_count"`
	FollowingCount int64   `json:"following_count"`
	TotalPostCount int64   `json:"total_post_count"`
	Posts          []Post  `json:"posts"`
}

type Post struct {
	ID          int64
	Content     *string  // text content of the post
	Attachments []string // list of images/video url's
	Likes       int64
	Comments    int64
	Reposts     int64
	Bookmarks   int64
	CreatedAt   time.Time
}

type SearchUserRequest struct {
	Rank    *string `json:"rank" binding:"omitempty,safe_alphabets"`
	Name    *string `json:"name" binding:"omitempty,safe_alphabets"`
	Company *string `json:"company" binding:"omitempty,alphanumspace"`
}

type SearchUserResponse struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Rank           string `json:"rank"`
	CompanyName    string `json:"company_name"`
}
