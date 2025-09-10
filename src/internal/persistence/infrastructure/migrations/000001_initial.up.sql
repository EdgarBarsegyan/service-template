-- BEGIN/COMMIT не нужен, так как используется утилита golag-migrate
CREATE SCHEMA IF NOT EXISTS migrations;

CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS
    app.users (
        id UUID NOT NULL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_users_username ON app.users (username);

COMMENT ON TABLE app.users IS 'Содержит информацию о пользователях системы';

COMMENT ON COLUMN app.users.id IS 'Идентификатор пользователя';

COMMENT ON COLUMN app.users.username IS 'Имя пользоватебя';

COMMENT ON COLUMN app.users.created_at IS 'Дата созодания';

COMMENT ON COLUMN app.users.updated_at IS 'Дата последнего обновления';