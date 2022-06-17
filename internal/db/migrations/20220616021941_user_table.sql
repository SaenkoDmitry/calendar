-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name  TEXT        NOT NULL,
    second_name TEXT        NOT NULL,
    email       TEXT UNIQUE NOT NULL,
    user_zone   TEXT        NOT NULL DEFAULT 'UTC'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
