-- Önce frozen kolonlarını kaldır
ALTER TABLE users
    DROP COLUMN IF EXISTS frozen_reason,
    DROP COLUMN IF EXISTS frozen_date;

-- Enum tiplerini kaldır
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role; 