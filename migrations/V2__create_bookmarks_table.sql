CREATE TABLE IF NOT EXISTS bookmarks (
                                         id BIGSERIAL PRIMARY KEY,
                                         user_id BIGINT NOT NULL,
                                         target_type VARCHAR(20) NOT NULL,
    target_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_bookmarks_user_target UNIQUE (user_id, target_type, target_id)
    );

CREATE INDEX IF NOT EXISTS idx_bookmarks_user
    ON bookmarks (user_id);