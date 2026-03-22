CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    rank VARCHAR(255),
    seman_book TEXT,
    company_name VARCHAR(255),
    licence_url TEXT,
    licence_type VARCHAR(255),
    licence_number VARCHAR(255),
    profile_picture TEXT,
    bio VARCHAR(1000),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- UPDATED AT TRIGGER FOR USERS
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ==========================================================

CREATE TABLE user_profile_summary (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    followers_count BIGINT NOT NULL DEFAULT 0,
    following_count BIGINT NOT NULL DEFAULT 0,
    total_post_count INTEGER DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION trg_init_user_profile_summary()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO user_profile_summary (user_id) VALUES (NEW.id);
    RETURN NULL;
END;
$$;

CREATE TRIGGER trg_user_profile_summary_init
AFTER INSERT ON users
FOR EACH ROW EXECUTE FUNCTION trg_init_user_profile_summary();

-- ==========================================================

CREATE TABLE user_follows (
    follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (follower_id, followee_id),

    CONSTRAINT no_self_follow CHECK (follower_id <> followee_id)
);

-- Speed up "who does user X follow?" queries
CREATE INDEX idx_user_follows_follower ON user_follows(follower_id);
-- Speed up "who follows user X?" queries
CREATE INDEX idx_user_follows_followee
ON user_follows(followee_id, follower_id);

-- Trigger to keep the follower and following count in sync
CREATE OR REPLACE FUNCTION sync_follow_counts()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Bump follower's following_count
        INSERT INTO user_profile_summary (user_id, following_count)
            VALUES (NEW.follower_id, 1)
            ON CONFLICT (user_id)
            DO UPDATE SET following_count = user_profile_summary.following_count + 1,
                          updated_at      = NOW();

        -- Bump followee's followers_count
        INSERT INTO user_profile_summary (user_id, followers_count)
            VALUES (NEW.followee_id, 1)
            ON CONFLICT (user_id)
            DO UPDATE SET followers_count = user_profile_summary.followers_count + 1,
                          updated_at      = NOW();

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE user_profile_summary
            SET following_count = GREATEST(following_count - 1, 0), updated_at = NOW()
            WHERE user_id = OLD.follower_id;

        UPDATE user_profile_summary
            SET followers_count = GREATEST(followers_count - 1, 0), updated_at = NOW()
            WHERE user_id = OLD.followee_id;
    END IF;

    RETURN NULL;
END;
$$;

CREATE TRIGGER trg_sync_follow_counts
AFTER INSERT OR DELETE ON user_follows
FOR EACH ROW EXECUTE FUNCTION sync_follow_counts();

