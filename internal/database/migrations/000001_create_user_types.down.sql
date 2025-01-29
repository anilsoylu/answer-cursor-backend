-- Users tablosunu eski haline getir
ALTER TABLE users 
    ALTER COLUMN status TYPE varchar(10) USING status::varchar(10),
    ALTER COLUMN role TYPE varchar(15) USING role::varchar(15);

-- Enum tipleri sil
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role; 