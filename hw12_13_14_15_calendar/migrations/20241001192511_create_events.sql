-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS events
(
    id            UUID PRIMARY KEY,
    title         TEXT      NOT NULL,
    datetime      TIMESTAMP NOT NULL,
    duration      INTERVAL  NOT NULL,
    description   TEXT,
    user_id       UUID      NOT NULL,
    notify_before INTERVAL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS events
-- +goose StatementEnd
