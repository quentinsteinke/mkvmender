-- MKV Mender Initial Schema

-- Users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    api_key TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_api_key ON users(api_key);

-- File hashes table
CREATE TABLE IF NOT EXISTS file_hashes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hash TEXT NOT NULL UNIQUE,
    file_size INTEGER NOT NULL,
    media_type TEXT NOT NULL CHECK(media_type IN ('movie', 'tv')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_file_hashes_hash ON file_hashes(hash);
CREATE INDEX IF NOT EXISTS idx_file_hashes_media_type ON file_hashes(media_type);

-- Naming submissions table
CREATE TABLE IF NOT EXISTS naming_submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hash_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    filename TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (hash_id) REFERENCES file_hashes(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_naming_submissions_hash_id ON naming_submissions(hash_id);
CREATE INDEX IF NOT EXISTS idx_naming_submissions_user_id ON naming_submissions(user_id);

-- Votes table
CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    submission_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    vote_type INTEGER NOT NULL CHECK(vote_type IN (-1, 1)),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(submission_id, user_id),
    FOREIGN KEY (submission_id) REFERENCES naming_submissions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_votes_submission_id ON votes(submission_id);
CREATE INDEX IF NOT EXISTS idx_votes_user_id ON votes(user_id);

-- Naming metadata table for additional information
CREATE TABLE IF NOT EXISTS naming_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    submission_id INTEGER NOT NULL UNIQUE,
    title TEXT,
    year INTEGER,
    season INTEGER,
    episode INTEGER,
    quality TEXT,
    source TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (submission_id) REFERENCES naming_submissions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_naming_metadata_submission_id ON naming_metadata(submission_id);

-- View for submissions with vote counts
CREATE VIEW IF NOT EXISTS submissions_with_votes AS
SELECT
    ns.id,
    ns.hash_id,
    ns.user_id,
    ns.filename,
    ns.created_at,
    fh.hash,
    fh.file_size,
    fh.media_type,
    u.username,
    COALESCE(SUM(v.vote_type), 0) as vote_score,
    COUNT(CASE WHEN v.vote_type = 1 THEN 1 END) as upvotes,
    COUNT(CASE WHEN v.vote_type = -1 THEN 1 END) as downvotes
FROM naming_submissions ns
JOIN file_hashes fh ON ns.hash_id = fh.id
JOIN users u ON ns.user_id = u.id
LEFT JOIN votes v ON ns.id = v.submission_id
GROUP BY ns.id, ns.hash_id, ns.user_id, ns.filename, ns.created_at,
         fh.hash, fh.file_size, fh.media_type, u.username;
