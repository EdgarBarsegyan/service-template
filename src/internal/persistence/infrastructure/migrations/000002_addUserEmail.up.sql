-- BEGIN/COMMIT не нужен, так как используется утилита golag-migrate
ALTER TABLE app.users 
ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE NULL;

COMMENT ON COLUMN app.users.email IS 'Электронная почта пользователя';