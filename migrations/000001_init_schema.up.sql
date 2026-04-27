-- 000001_init_schema.up.sql

CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_like_unique ON likes (user_id, target_type, target_id);

CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bookmark_unique ON bookmarks (user_id, target_type, target_id);

CREATE TABLE IF NOT EXISTS reports (
    id SERIAL PRIMARY KEY,
    reporter_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    target_author_id INTEGER NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
