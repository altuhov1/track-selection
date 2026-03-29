-- +goose Up
-- +goose StatementBegin
CREATE TABLE track_teachers (
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (track_id, teacher_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS track_teachers CASCADE;
-- +goose StatementEnd