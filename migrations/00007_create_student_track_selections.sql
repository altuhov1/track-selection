-- +goose Up
CREATE TABLE IF NOT EXISTS student_track_selections (
    id VARCHAR(36) PRIMARY KEY,
    student_id VARCHAR(36) NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    track_id VARCHAR(36) NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    selected_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(student_id, track_id)
);

CREATE INDEX idx_student_track_selections_student_id ON student_track_selections(student_id);
CREATE INDEX idx_student_track_selections_track_id ON student_track_selections(track_id);

COMMENT ON TABLE student_track_selections IS 'Выбранные студентом треки';
COMMENT ON COLUMN student_track_selections.student_id IS 'ID студента';
COMMENT ON COLUMN student_track_selections.track_id IS 'ID выбранного трека';
COMMENT ON COLUMN student_track_selections.selected_at IS 'Дата выбора';

-- +goose Down
DROP TABLE IF EXISTS student_track_selections;