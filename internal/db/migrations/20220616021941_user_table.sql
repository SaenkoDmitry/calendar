-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id    SERIAL PRIMARY KEY,
    username  TEXT,
    email TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users IF EXISTS;
-- +goose StatementEnd
