-- +goose Up
DROP TABLE IF EXISTS tracks;

CREATE TABLE tracks (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    
    -- Учебный план (JSON с семестрами и курсами)
    curriculum JSONB NOT NULL DEFAULT '{"semesters": []}',
    
    -- Требования
    requirements JSONB NOT NULL DEFAULT '[]',
    
    -- Преподаватели
    teachers JSONB NOT NULL DEFAULT '[]',
    
    -- Сложность (1-5)
    difficulty INT NOT NULL DEFAULT 3,
    
    -- Тип трека (для фронтенда)
    type INT NOT NULL DEFAULT 0,
    
    employment_prospects INT NOT NULL DEFAULT 5,
    alumni_reviews INT NOT NULL DEFAULT 5,
    web_link VARCHAR(512),
    has_certificates INT NOT NULL DEFAULT 0,
    learning_style INT NOT NULL DEFAULT 1,
    desired_tech_skills INT NOT NULL DEFAULT 5,
    desired_math_skills INT NOT NULL DEFAULT 5,
    desired_soft_skills INT NOT NULL DEFAULT 5,
    professional_goals JSONB NOT NULL DEFAULT '[]',
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Добавляем CHECK constraints отдельно (чтобы избежать проблем с DROP)
ALTER TABLE tracks ADD CONSTRAINT tracks_difficulty_check CHECK (difficulty BETWEEN 1 AND 5);
ALTER TABLE tracks ADD CONSTRAINT tracks_employment_prospects_check CHECK (employment_prospects BETWEEN 1 AND 10);
ALTER TABLE tracks ADD CONSTRAINT tracks_alumni_reviews_check CHECK (alumni_reviews BETWEEN 1 AND 10);
ALTER TABLE tracks ADD CONSTRAINT tracks_has_certificates_check CHECK (has_certificates IN (0, 1));
ALTER TABLE tracks ADD CONSTRAINT tracks_learning_style_check CHECK (learning_style IN (1, 2, 3));
ALTER TABLE tracks ADD CONSTRAINT tracks_desired_tech_skills_check CHECK (desired_tech_skills BETWEEN 1 AND 10);
ALTER TABLE tracks ADD CONSTRAINT tracks_desired_math_skills_check CHECK (desired_math_skills BETWEEN 1 AND 10);
ALTER TABLE tracks ADD CONSTRAINT tracks_desired_soft_skills_check CHECK (desired_soft_skills BETWEEN 1 AND 10);

-- Создаем индексы
CREATE INDEX idx_tracks_name ON tracks(name);
CREATE INDEX idx_tracks_difficulty ON tracks(difficulty);
CREATE INDEX idx_tracks_type ON tracks(type);

-- +goose Down
DROP TABLE IF EXISTS tracks;