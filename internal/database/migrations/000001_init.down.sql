-- Drop trigger on user_follows
DROP TRIGGER IF EXISTS trg_sync_follow_counts ON user_follows;

-- Drop function for follow counts
DROP FUNCTION IF EXISTS sync_follow_counts;

-- Drop indexes
DROP INDEX IF EXISTS idx_user_follows_followee;
DROP INDEX IF EXISTS idx_user_follows_follower;

-- Drop user_follows table
DROP TABLE IF EXISTS user_follows;

-- Drop user_profile_summary table
DROP TABLE IF EXISTS user_profile_summary;

-- Drop trigger on users
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function for updated_at
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop users table
DROP TABLE IF EXISTS users;