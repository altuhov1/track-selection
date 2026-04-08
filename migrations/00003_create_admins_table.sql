-- +goose Up
CREATE TABLE IF NOT EXISTS admins (
    id VARCHAR(36) PRIMARY KEY,
    auth_user_id VARCHAR(36) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_admins_auth_user FOREIGN KEY (auth_user_id) REFERENCES auth_users(id) ON DELETE CASCADE
);

CREATE INDEX idx_admins_auth_user_id ON admins(auth_user_id);
CREATE INDEX idx_admins_email ON admins(email);

COMMENT ON TABLE admins IS 'Профили администраторов';
COMMENT ON COLUMN admins.id IS 'Уникальный идентификатор администратора';
COMMENT ON COLUMN admins.auth_user_id IS 'Ссылка на учетную запись (auth_users)';
COMMENT ON COLUMN admins.email IS 'Email администратора';
COMMENT ON COLUMN admins.first_name IS 'Имя администратора';
COMMENT ON COLUMN admins.last_name IS 'Фамилия администратора';
COMMENT ON COLUMN admins.created_at IS 'Дата создания профиля';
COMMENT ON COLUMN admins.updated_at IS 'Дата последнего обновления';

-- +goose Down
DROP TABLE IF EXISTS admins;