-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_meetings DROP CONSTRAINT user_meetings_pkey;
ALTER TABLE user_meetings DROP COLUMN id;
ALTER TABLE user_meetings ADD PRIMARY KEY (user_id, meeting_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_meetings DROP CONSTRAINT user_meetings_pkey;
ALTER TABLE user_meetings ADD COLUMN id INT GENERATED ALWAYS AS IDENTITY;
ALTER TABLE user_meetings ADD PRIMARY KEY (id);
-- +goose StatementEnd
