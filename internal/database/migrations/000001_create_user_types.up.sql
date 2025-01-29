-- Kullanıcı durumu için enum tipi
DO $$ BEGIN
    DROP TYPE IF EXISTS user_status;
    CREATE TYPE user_status AS ENUM ('active', 'passive', 'banned', 'frozen');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Kullanıcı rolü için enum tipi
DO $$ BEGIN
    DROP TYPE IF EXISTS user_role;
    CREATE TYPE user_role AS ENUM ('USER', 'EDITOR', 'ADMIN', 'SUPER_ADMIN');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Users tablosunu güncelle
ALTER TABLE users 
    ALTER COLUMN status TYPE user_status USING status::text::user_status,
    ALTER COLUMN role TYPE user_role USING role::text::user_role,
    ADD COLUMN IF NOT EXISTS frozen_reason TEXT,
    ADD COLUMN IF NOT EXISTS frozen_date TIMESTAMP; 