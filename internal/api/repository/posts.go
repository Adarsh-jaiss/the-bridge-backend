package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/the-bridge/internal/api/types/models"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type PostRepository interface {
	CreatePost(ctx context.Context, userID int64, post models.Post) error
	SoftDeletePost(ctx context.Context, postID, userID int64) error
	ToggleLike(ctx context.Context, postID, userID int64) (bool, error)
	LikePost(ctx context.Context, postID, userID int64) error
	UnlikePost(ctx context.Context, postID, userID int64) error
	RepostPost(ctx context.Context, postId, repost_account_id int64) error
	BookmarkPost(ctx context.Context, postId, bookmark_account_id int64) error
	ReportPost(ctx context.Context, reportedBy, postId int64, reason string) error

	// Todo
	AddComments() // single level
	DeleteComments()
	LikeComment()
	UnlikeComment()
	SelectPollOptions()
	FetchPollMetrics()
	FetchUserFeed()
	CreateNotification()
	FetchPostStats() // likes,comments,bookmarks, reposts
	CreateAnnouncements() // for all users -> by admin
	FetchAnnouncememts() // for all users -> by admin
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{
		db: db,
	}
}

func (p *postRepository) CreatePost(ctx context.Context, userID int64, post models.Post) error {
	log := logger.FromContext(ctx)

	insertquery := `INSERT INTO posts (
		user_id,
		post_type,
		content,
		attachments
	) VALUES ($1,$2,$3,$4)
	RETURNING id;
	`
	pollQuery := `INSERT INTO poll_options (
		post_id, 
		option_text, 
		position
	) SELECT $1, unnest($2::text[]), 
	generate_series(1, array_length($2,1));`

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("error starting a transcation")
		return err
	}
	defer tx.Rollback().Error()

	var postId int64
	err = tx.QueryRowContext(ctx, insertquery,
		userID,
		post.Type,
		post.Content,
		pq.Array(post.Attachments),
	).Scan(&postId)
	if err != nil {
		log.Error("error creating post", zap.String("err", err.Error()))
		return err
	}

	if post.Type == "poll" && len(post.PollOptions) > 0 {
		_, err := tx.ExecContext(ctx, pollQuery, postId, pq.Array(post.PollOptions))

		if err != nil {
			log.Error("failed to insert poll options", zap.Error(err))
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (p *postRepository) SoftDeletePost(ctx context.Context, postID, userID int64) error {
	log := logger.FromContext(ctx)

	query := `
	UPDATE posts
	SET is_deleted = TRUE,
	    deleted_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND is_deleted = FALSE;
	`

	res, err := p.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		log.Error("failed to soft delete post", zap.Error(err))
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Error("err", zap.Error(err))
		return err
	}

	if rows == 0 {
		log.Error("err", zap.Error(err))
		return fmt.Errorf("post not found or not owned by user")
	}
	return nil
}

func (p *postRepository) LikePost(ctx context.Context, postID, userID int64) error {
	log := logger.FromContext(ctx)

	query := `
	INSERT INTO post_likes (post_id, user_id)
	VALUES ($1, $2)
	ON CONFLICT (post_id, user_id) DO NOTHING;
	`

	_, err := p.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		log.Error("failed to like post", zap.Error(err))
		return err
	}

	return nil
}

func (p *postRepository) UnlikePost(ctx context.Context, postID, userID int64) error {
	log := logger.FromContext(ctx)

	query := `
	DELETE FROM post_likes
	WHERE post_id = $1 AND user_id = $2;
	`

	_, err := p.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		log.Error("failed to unlike post", zap.Error(err))
		return err
	}

	return nil
}

func (p *postRepository) ToggleLike(ctx context.Context, postID, userID int64) (bool, error) {
	log := logger.FromContext(ctx)

	query := `
	WITH deleted AS (
		DELETE FROM post_likes
		WHERE post_id = $1 AND user_id = $2
		RETURNING 1
	)
	INSERT INTO post_likes (post_id, user_id)
	SELECT $1, $2
	WHERE NOT EXISTS (SELECT 1 FROM deleted)
	RETURNING 1;
	`

	var liked int
	err := p.db.QueryRowContext(ctx, query, postID, userID).Scan(&liked)

	if err == sql.ErrNoRows {
		return false, nil // unliked
	}
	if err != nil {
		log.Error("failed to toggle like", zap.Error(err))
		return false, err
	}

	return true, nil // liked
}

func (p *postRepository) RepostPost(ctx context.Context, postId, repost_account_id int64) error {
	log := logger.FromContext(ctx)

	query := `INSERT INTO post_reposts (
		post_id,
		user_id
	)VALUES($1,$2);`

	_, err := p.db.ExecContext(ctx, query, postId, repost_account_id)
	if err != nil {
		log.Error("error creating repost", zap.Error(err))
		return err
	}

	return nil
}

func (p *postRepository) BookmarkPost(ctx context.Context, postId, bookmark_account_id int64) error {
	log := logger.FromContext(ctx)

	query := `INSERT INTO post_bookmarks (
		post_id,
		user_id
	)VALUES($1,$2);`

	_, err := p.db.ExecContext(ctx, query, postId, bookmark_account_id)
	if err != nil {
		log.Error("error creating repost", zap.Error(err))
		return err
	}

	return nil
}

func (p *postRepository) ReportPost(ctx context.Context, reportedBy, postId int64, reason string) error {
	log := logger.FromContext(ctx)
	updatePostsQuery := `UPDATE posts SET is_reported = TRUE AND updated_at = NOW();`

	insertReportQuery := `INSERT INTO post_reports (
		post_id,
		reported_by,
		reason
	)VALUES($1, $2, $3)
	`
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("error starting a transcation")
		return err
	}
	defer tx.Rollback().Error()

	_, err = tx.ExecContext(ctx, updatePostsQuery, nil)
	if err != nil {
		log.Error("failed to update post report status", zap.Error(err))
	}

	_, err = tx.ExecContext(ctx, insertReportQuery, postId, reportedBy, reason)
	if err != nil {
		log.Error("failed to update post report status", zap.Error(err))
	}

	if err = tx.Commit(); err != nil {
		log.Error("err", zap.Error(err))
		return err
	}

	return nil

}
