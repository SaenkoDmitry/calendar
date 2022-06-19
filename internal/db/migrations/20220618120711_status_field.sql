-- +goose Up
-- +goose StatementBegin
CREATE TYPE status_enum AS ENUM ('requested', 'approved', 'declined');
ALTER TABLE user_meetings
    ADD status status_enum NOT NULL DEFAULT 'requested';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_meetings DROP COLUMN status;
DROP TYPE status_enum;
-- +goose StatementEnd
