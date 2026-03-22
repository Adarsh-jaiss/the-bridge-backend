
-- TRIGGERS
DROP TRIGGER IF EXISTS trg_repost_count      ON post_reposts;
DROP TRIGGER IF EXISTS trg_like_count        ON post_likes;
DROP TRIGGER IF EXISTS trg_comments_count    ON comments;
DROP TRIGGER IF EXISTS trg_posts_stats_init  ON posts;


-- TRIGGER FUNCTIONS
DROP FUNCTION IF EXISTS trg_sync_repost_count();
DROP FUNCTION IF EXISTS trg_sync_like_count();
DROP FUNCTION IF EXISTS trg_sync_comments_count();
DROP FUNCTION IF EXISTS trg_init_posts_stats();


-- INDEXES
DROP INDEX IF EXISTS idx_bookmarks_user_id;
DROP INDEX IF EXISTS idx_post_reports_reportedby;
DROP INDEX IF EXISTS idx_post_reports_status;
DROP INDEX IF EXISTS idx_post_reports_post_id;
DROP INDEX IF EXISTS idx_comments_parent_id;
DROP INDEX IF EXISTS idx_comments_post_id;
DROP INDEX IF EXISTS idx_posts_user_created;
DROP INDEX IF EXISTS idx_posts_feed;


-- TABLES (child tables first, then parents)
DROP TABLE IF EXISTS post_bookmarks;
DROP TABLE IF EXISTS post_reposts;
DROP TABLE IF EXISTS post_likes;
DROP TABLE IF EXISTS post_reports;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts_stats;
DROP TABLE IF EXISTS posts;