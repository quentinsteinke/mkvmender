-- MKV Mender Admin Features Migration

-- Add role column to users table
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user' CHECK(role IN ('user', 'moderator', 'admin'));

-- Add is_active column for user account suspension
ALTER TABLE users ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT 1;

-- Create index on role for quick admin queries
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- Moderation actions table for audit trail
CREATE TABLE IF NOT EXISTS moderation_actions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    admin_id INTEGER NOT NULL,
    action_type TEXT NOT NULL CHECK(action_type IN (
        'delete_submission',
        'change_role',
        'suspend_user',
        'activate_user'
    )),
    target_type TEXT NOT NULL CHECK(target_type IN ('user', 'submission')),
    target_id INTEGER NOT NULL,
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_moderation_actions_admin_id ON moderation_actions(admin_id);
CREATE INDEX IF NOT EXISTS idx_moderation_actions_target ON moderation_actions(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_moderation_actions_created_at ON moderation_actions(created_at DESC);
