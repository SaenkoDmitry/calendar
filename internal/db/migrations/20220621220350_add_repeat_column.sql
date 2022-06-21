-- +goose Up
-- +goose StatementBegin
CREATE TYPE repeat_enum AS ENUM ('days', 'weeks', 'months', 'years', 'weekdays');
ALTER TABLE meetings
    ADD repeat repeat_enum;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meetings DROP COLUMN repeat;
DROP TYPE repeat_enum;
-- +goose StatementEnd
