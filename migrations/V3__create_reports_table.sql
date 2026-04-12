CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    reporter_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    target_author_id INTEGER NOT NULL,
    target_type VARCHAR(50) NOT NULL, -- 'article', 'post', 'comment'
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'OPEN', -- 'OPEN', 'REJECTED', 'RESOLVED'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_target ON reports(target_id, target_type);
