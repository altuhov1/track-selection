-- +goose Up
CREATE TABLE IF NOT EXISTS profile_completion (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE REFERENCES auth_users(id) ON DELETE CASCADE,
    is_complete BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_profile_completion_user_id ON profile_completion(user_id);
CREATE INDEX idx_profile_completion_is_complete ON profile_completion(is_complete);

COMMENT ON TABLE profile_completion IS 'Статус заполнения профиля пользователя';
COMMENT ON COLUMN profile_completion.is_complete IS 'TRUE - все данные заполнены, FALSE - не все';
COMMENT ON COLUMN profile_completion.completed_at IS 'Дата, когда профиль стал полностью заполненным';

-- +goose Down
DROP TABLE IF EXISTS profile_completion;