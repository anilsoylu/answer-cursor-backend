-- Drop existing unique constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
DROP INDEX IF EXISTS idx_username_deleted;
DROP INDEX IF EXISTS idx_email_deleted;

-- Create new unique constraints that ignore soft deleted and frozen records
CREATE UNIQUE INDEX idx_username_active ON users (username) 
WHERE deleted_at IS NULL AND status != 'frozen';

CREATE UNIQUE INDEX idx_email_active ON users (email) 
WHERE deleted_at IS NULL AND status != 'frozen'; 