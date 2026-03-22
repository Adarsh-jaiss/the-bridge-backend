package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/adarsh-jaiss/the-bridge/internal/api/types/models"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/lib/pq"
	pg "github.com/lib/pq"
	"go.uber.org/zap"
)

type UserProfile interface {
	CreateProfile(ctx context.Context, userId int64, user models.User) error
	Follow(ctx context.Context, followeeID, followerID int64) error
	UnFollow(ctx context.Context, followeeID, followerID int64) error
	UpdateBioAndProfilePic(ctx context.Context, bio, profilePicture string, userID int64) error
	FetchUserProfileSummary(ctx context.Context, userID int64) (*models.UserProfile, error)
	FetchUserPosts(ctx context.Context, limit, cursor, userID int64) ([]models.Post, int64, error) // int64 is cursor
}

type User struct {
	db *sql.DB
}

func NewUserProfile(db *sql.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) CreateProfile(ctx context.Context, userId int64, user models.User) error {
	log := logger.FromContext(ctx)

	query := `
		UPDATE users SET
			first_name = $1,
			last_name = $2,
			rank = $3,
			seman_book = $4,
			company_name = $5,
			licence_url = $6,
			licence_type = $7,
			licence_number = $8
		WHERE id = $9
	`

	result, err := u.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Rank,
		user.SemanBook,
		user.CompanyName,
		user.Licence,
		user.LicenceType,
		user.LicenceNumber,
		userId)

	if err != nil {
		log.Error("error inserting user data in db", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("error checking rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		log.Error("couldn't updated db for create user, user not found!", zap.Any("user_id", userId))
		return fmt.Errorf("user not found")
	}

	return nil

}

func (u *User) Follow(ctx context.Context, followeeID, followerID int64) error {
	log := logger.FromContext(ctx)

	query := `INSERT INTO user_follows(follower_id, followee_id) VALUES ($1,$2)`

	_, err := u.db.ExecContext(ctx, query, followerID, followeeID)
	if err != nil {
		if pqErr, ok := err.(*pg.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation — already following
				return fmt.Errorf("already following user")
			case "23503": // foreign_key_violation — user doesn't exist
				return fmt.Errorf("user not found")
			}
		}
		log.Error("error inserting user follow", zap.Int64("follower_id", followerID),
			zap.Int64("followee_id", followerID), zap.Error(err))
		return err
	}

	return nil

}

func (u *User) UnFollow(ctx context.Context, followeeID, followerID int64) error {
	log := logger.FromContext(ctx)

	query := `DELETE FROM user_follows WHERE follower_id = $1 AND followee_id = $2`

	res, err := u.db.ExecContext(ctx, query, followerID, followeeID)
	if err != nil {
		log.Error("error deleting user follow",
			zap.Int64("follower_id", followerID),
			zap.Int64("followee_id", followeeID),
			zap.Error(err))
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error("error checking rows affected", zap.Error(err))
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("not following user")
	}

	return nil
}

func (u *User) UpdateBioAndProfilePic(ctx context.Context, bio, profilePicture string, userID int64) error {
	log := logger.FromContext(ctx)
	query := `UPDATE users SET `

	args := []any{}
	i := 1

	if bio != "" {
		query += fmt.Sprintf("bio=$%d", i)
		args = append(args, bio)
		i++
	}
	if profilePicture != "" {
		query += fmt.Sprintf("profile_picture=$%d", i)
		args = append(args, profilePicture)
		i++
	}

	if len(args) == 0 {
		return nil //nothing to update
	}

	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id=$%d", i)
	args = append(args, userID)

	res, err := u.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Error("error updating user profile", zap.Error(err))
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error("error checking rows affected", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		log.Error("Error updating profile,user not found", zap.Int64("user_id", userID))
	}
	return nil
}

func (u *User) FetchUserProfileSummary(ctx context.Context, userID int64) (*models.UserProfile, error) {
	log := logger.FromContext(ctx)
	query := `SELECT u.first_name, 
		u.last_name, 
		u.profile_picture,
		u.rank, 
		u.bio,
		s.followers_count,
		s.following_count,
		s.total_post_count
	FROM users AS u
	INNER JOIN user_profile_summary AS s
		ON u.id=s.user_id
	WHERE u.id=$1;	
	`

	var user models.UserProfile
	if err := u.db.QueryRowContext(ctx, query, userID).Scan(
		&user.FirstName,
		&user.LastName,
		&user.ProfilePicture,
		&user.Rank,
		&user.Bio,
		&user.FollowersCount,
		&user.FollowingCount,
		&user.TotalPostCount,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Warn("user dosen't exists",zap.Int64("user_id",userID))
			return nil, fmt.Errorf("user don't exists")
		}
		log.Error("error fetching user profile summary", zap.Int64("user_id", userID), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (u *User) FetchUserPosts(ctx context.Context, limit, cursor, userID int64) ([]models.Post, int64, error) {
	log := logger.FromContext(ctx)
	query := `SELECT p.id,
		p.content,
		p.attachments,
		p.created_at,
		s.likes_count,
		s.comments_count,
		s.reposts_count,
		s.bookmarks_count
	FROM posts AS p
	INNER JOIN posts_stats AS s
		ON p.id=s.post_id
	WHERE p.user_id = $1 
		AND p.is_deleted = FALSE
		AND ($2 = 0 OR p.id < $2)  
	ORDER BY p.id DESC
	LIMIT $3
	`
	rows, err := u.db.QueryContext(ctx, query, userID, cursor, limit)
	if err != nil {
		log.Error("failed to fetch user posts", zap.Error(err), zap.Int64("user_id", userID), zap.Int64("cursor", cursor))
		return nil, 0, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&p.ID,
			&p.Content,
			pq.Array(&p.Attachments), // ← wraps []string so pq handles the conversion
			&p.CreatedAt,
			&p.Likes,
			&p.Comments,
			&p.Reposts,
			&p.Bookmarks,
		); err != nil {
			log.Error("failed to scan post row", zap.Error(err))
			return nil, 0, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		log.Error("row iteration error", zap.Error(err))
		return nil, 0, err
	}

	var nextCursor int64
	if len(posts) == int(limit) {
		nextCursor = posts[len(posts)-1].ID
	}
	return posts, nextCursor, nil

}
