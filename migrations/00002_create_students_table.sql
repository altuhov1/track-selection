-- +goose Up
CREATE TABLE IF NOT EXISTS students (
    id VARCHAR(36) PRIMARY KEY,
    auth_user_id VARCHAR(36) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    rating INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_students_auth_user FOREIGN KEY (auth_user_id) REFERENCES auth_users(id) ON DELETE CASCADE
);

CREATE INDEX idx_students_auth_user_id ON students(auth_user_id);
CREATE INDEX idx_students_email ON students(email);
CREATE INDEX idx_students_username ON students(username);

COMMENT ON TABLE students IS 'Профили студентов';
COMMENT ON COLUMN students.id IS 'Уникальный идентификатор студента';
COMMENT ON COLUMN students.auth_user_id IS 'Ссылка на учетную запись (auth_users)';
COMMENT ON COLUMN students.email IS 'Email студента';
COMMENT ON COLUMN students.first_name IS 'Имя студента';
COMMENT ON COLUMN students.last_name IS 'Фамилия студента';
COMMENT ON COLUMN students.username IS 'Имя пользователя (отображается в интерфейсе)';
COMMENT ON COLUMN students.rating IS 'Рейтинг студента (для выбора треков, 0-100)';
COMMENT ON COLUMN students.created_at IS 'Дата создания профиля';
COMMENT ON COLUMN students.updated_at IS 'Дата последнего обновления';

-- +goose Down
DROP TABLE IF EXISTS students;