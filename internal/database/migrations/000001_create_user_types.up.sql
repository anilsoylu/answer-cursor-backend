-- Kullanıcı durumu için enum tipi
DO $$ BEGIN
    CREATE TYPE user_status AS ENUM ('active', 'passive', 'banned');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Kullanıcı rolü için enum tipi
DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('USER', 'EDITOR', 'ADMIN', 'SUPER_ADMIN');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Users tablosunu güncelle
ALTER TABLE users 
    ALTER COLUMN status TYPE user_status USING status::user_status,
    ALTER COLUMN role TYPE user_role USING role::user_role; 