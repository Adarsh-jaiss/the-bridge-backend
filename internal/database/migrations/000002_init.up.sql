CREATE TABLE posts (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT,
    is_reported BOOLEAN DEFAULT FALSE,
    attachments TEXT[],
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted  BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT chk_post_not_empty CHECK (
        content IS NOT NULL OR (attachments IS NOT NULL AND array_length(attachments, 1) > 0)
    )
);


CREATE INDEX idx_posts_feed
    ON posts(created_at DESC)
    WHERE is_deleted = FALSE;             -- global/home feed ordered by recency

-- composite for "a specific user's posts ordered by time" (profile feed)
CREATE INDEX idx_posts_user_created
    ON posts(user_id, created_at DESC)
    WHERE is_deleted = FALSE;

-- ============================================================

CREATE TABLE posts_stats (
    post_id BIGINT PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    likes_count BIGINT NOT NULL DEFAULT 0,
    comments_count BIGINT NOT NULL DEFAULT 0,
    reposts_count BIGINT NOT NULL DEFAULT 0,
    bookmarks_count BIGINT NOT NULL DEFAULT 0 
);

-- AUTO-CREATE posts_stats ROW WHEN A POST IS INSERTED
-- Required so all the count triggers below have a row to update

CREATE OR REPLACE FUNCTION trg_init_posts_stats()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO posts_stats (post_id) VALUES (NEW.id) ON CONFLICT (post_id) DO NOTHING;
    RETURN NULL;
END;
$$;

CREATE TRIGGER trg_posts_stats_init
AFTER INSERT ON posts
FOR EACH ROW EXECUTE FUNCTION trg_init_posts_stats();

-- ============================================================

CREATE TABLE post_reports (
    report_id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    post_id     BIGINT REFERENCES posts(id) ON DELETE CASCADE,
    reported_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reason      TEXT,
    status      TEXT NOT NULL DEFAULT 'pending'   -- added: pending | reviewed | dismissed | actioned
                CHECK (status IN ('pending', 'reviewed', 'dismissed', 'actioned')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),  -- added: when was it reported

    CONSTRAINT uq_report_per_user UNIQUE (post_id, reported_by)  -- one report per user per post
);

-- post_reports: moderation dashboard queries
CREATE INDEX idx_post_reports_post_id    ON post_reports(post_id);
CREATE INDEX idx_post_reports_status     ON post_reports(status)     -- "show all pending reports"
    WHERE status = 'pending';
CREATE INDEX idx_post_reports_reportedby ON post_reports(reported_by);



-- ============================================================

CREATE TABLE comments (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    post_id    BIGINT REFERENCES posts(id)    ON DELETE CASCADE,
    user_id    BIGINT REFERENCES users(id)    ON DELETE CASCADE,
    parent_id  BIGINT REFERENCES comments(id) ON DELETE CASCADE, -- NULL = top-level; set for replies
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ
);


CREATE INDEX idx_comments_post_id   ON comments(post_id)   WHERE is_deleted = FALSE;
CREATE INDEX idx_comments_parent_id ON comments(parent_id) WHERE is_deleted = FALSE;

-- NOTE: trigger only fires on hard DELETE, not soft delete
-- For soft delete, decrement the count in your application layer
-- when you SET is_deleted = TRUE, or use an UPDATE trigger (see below)
CREATE OR REPLACE FUNCTION trg_sync_comments_count()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.is_deleted = FALSE THEN
        UPDATE posts_stats SET comments_count = comments_count + 1 WHERE post_id = NEW.post_id;

    ELSIF TG_OP = 'UPDATE' THEN
        -- user just soft-deleted a comment
        IF OLD.is_deleted = FALSE AND NEW.is_deleted = TRUE THEN
            UPDATE posts_stats SET comments_count = comments_count - 1 WHERE post_id = NEW.post_id;
        -- user restored a soft-deleted comment
        ELSIF OLD.is_deleted = TRUE AND NEW.is_deleted = FALSE THEN
            UPDATE posts_stats SET comments_count = comments_count + 1 WHERE post_id = NEW.post_id;
        END IF;

    ELSIF TG_OP = 'DELETE' THEN
        -- only decrement if it wasn't already soft-deleted (avoid double decrement)
        IF OLD.is_deleted = FALSE THEN
            UPDATE posts_stats SET comments_count = comments_count - 1 WHERE post_id = OLD.post_id;
        END IF;
    END IF;
    RETURN NULL;
END;
$$;

CREATE TRIGGER trg_comments_count
AFTER INSERT OR UPDATE OF is_deleted OR DELETE ON comments
FOR EACH ROW EXECUTE FUNCTION trg_sync_comments_count();



-- ============================================================ 

CREATE TABLE post_likes (
    post_id    BIGINT NOT NULL REFERENCES posts(id)  ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)           -- one like per user per post, naturally de-duped
);

-- Trigger: keep posts_stats.like_count in sync
CREATE OR REPLACE FUNCTION trg_sync_like_count()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE posts_stats SET like_count = like_count + 1 WHERE post_id = NEW.post_id;

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE posts_stats SET like_count = like_count - 1 WHERE post_id = OLD.post_id;
    END IF;
    RETURN NULL;
END;
$$;

CREATE TRIGGER trg_like_count
AFTER INSERT OR DELETE ON post_likes
FOR EACH ROW EXECUTE FUNCTION trg_sync_like_count();

-- ============================================================ 

CREATE TABLE post_bookmarks (
    post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id) 
);

CREATE INDEX idx_bookmarks_user_id ON post_bookmarks(user_id);

-- ============================================================ 


CREATE TABLE post_reposts (
    post_id    BIGINT      NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);

CREATE OR REPLACE FUNCTION trg_sync_repost_count()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE posts_stats SET reposts_count = reposts_count + 1 WHERE post_id = NEW.post_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE posts_stats SET reposts_count = reposts_count - 1 WHERE post_id = OLD.post_id;
    END IF;
    RETURN NULL;
END;
$$;

CREATE TRIGGER trg_repost_count
AFTER INSERT OR DELETE ON post_reposts
FOR EACH ROW EXECUTE FUNCTION trg_sync_repost_count();



