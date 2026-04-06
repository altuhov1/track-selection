-- +goose Up
CREATE TABLE IF NOT EXISTS auth_users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('student', 'admin')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_users_email ON auth_users(email);

CREATE INDEX IF NOT EXISTS idx_auth_users_role ON auth_users(role);

COMMENT ON TABLE auth_users IS 'Таблица для хранения учетных записей пользователей';
COMMENT ON COLUMN auth_users.id IS 'Уникальный идентификатор пользователя (UUID)';
COMMENT ON COLUMN auth_users.email IS 'Email пользователя (уникальный)';
COMMENT ON COLUMN auth_users.password_hash IS 'Хеш пароля (bcrypt)';
COMMENT ON COLUMN auth_users.role IS 'Роль пользователя: student или admin';
COMMENT ON COLUMN auth_users.created_at IS 'Дата и время создания записи';
COMMENT ON COLUMN auth_users.updated_at IS 'Дата и время последнего обновления';

-- +goose Down
DROP TABLE IF EXISTS auth_users;