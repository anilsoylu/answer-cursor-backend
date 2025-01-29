-- Drop new unique indexes
DROP INDEX IF EXISTS idx_username_active;
DROP INDEX IF EXISTS idx_email_active;

-- Restore original unique constraints
ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email); 