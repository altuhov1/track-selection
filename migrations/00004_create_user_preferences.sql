-- +goose Up
CREATE TABLE IF NOT EXISTS user_preferences (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE REFERENCES auth_users(id) ON DELETE CASCADE,
    
    -- Профессиональные цели
    professional_goals JSONB DEFAULT '[]',
    
    -- Академические оценки (2-5)
    grades_informatics INT DEFAULT 0,
    grades_programming INT DEFAULT 0,
    grades_foreign_language INT DEFAULT 0,
    grades_physics INT DEFAULT 0,
    grades_aig INT DEFAULT 0,
    grades_math_analysis INT DEFAULT 0,
    grades_algorithms_data_structures INT DEFAULT 0,
    grades_databases INT DEFAULT 0,
    grades_discrete_math INT DEFAULT 0,
    
    -- Навыки (0-10)
    skills_databases INT DEFAULT 0,
    skills_system_architecture INT DEFAULT 0,
    skills_algorithmic_programming INT DEFAULT 0,
    skills_public_speaking INT DEFAULT 0,
    skills_testing INT DEFAULT 0,
    skills_analytics INT DEFAULT 0,
    skills_machine_learning INT DEFAULT 0,
    skills_os_knowledge INT DEFAULT 0,
    skills_research_projects INT DEFAULT 0,
    
    -- Стиль обучения (1-3)
    learning_style INT DEFAULT 1 CHECK (learning_style IN (1, 2, 3)),
    
    -- Сертификаты (0-1)
    certificates INT DEFAULT 0 CHECK (certificates IN (0, 1)),
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);

COMMENT ON TABLE user_preferences IS 'Предпочтения пользователя для подбора треков';
COMMENT ON COLUMN user_preferences.professional_goals IS 'Массив ID профессиональных целей';
COMMENT ON COLUMN user_preferences.learning_style IS '1-практика, 2-теория, 3-проектная работа';
COMMENT ON COLUMN user_preferences.certificates IS '0-нет, 1-да';

-- +goose Down
DROP TABLE IF EXISTS user_preferences;