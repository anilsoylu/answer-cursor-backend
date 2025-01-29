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
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255) DEFAULT '/uploads/default/avatar.png',
    status user_status NOT NULL,
    role user_role NOT NULL,
    is_root_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_date TIMESTAMP,
    ban_reason TEXT,
    ban_end_date TIMESTAMP,
    frozen_reason TEXT,
    frozen_date TIMESTAMP,
    deleted_at TIMESTAMP
); 