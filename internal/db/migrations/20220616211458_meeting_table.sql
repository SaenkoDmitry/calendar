-- +goose Up
-- +goose StatementBegin

CREATE TABLE meetings
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    meet_name   TEXT NOT NULL,
    description TEXT,
    start_date  DATE,
    start_time  TIME,
    end_date    DATE,
    end_time    TIME
);

CREATE TABLE user_meetings
(
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    INT NOT NULL,
    meeting_id INT NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_meeting_id FOREIGN KEY (meeting_id) REFERENCES meetings (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meetings;
DROP TABLE IF EXISTS user_meetings;
-- +goose StatementEnd
