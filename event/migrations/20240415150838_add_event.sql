-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event
(
    id          bigserial PRIMARY KEY,
    owner       uuid,
    members     uuid[],
    title       varchar(64),
    description text,
    date        timestamp
);

CREATE INDEX IF NOT EXISTS idx_event_owner ON event (owner);
CREATE INDEX IF NOT EXISTS idx_event_title ON event (title);
CREATE INDEX IF NOT EXISTS idx_event_date ON event (date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;

DROP INDEX IF EXISTS idx_event_owner;
DROP INDEX IF EXISTS idx_event_title;
DROP INDEX IF EXISTS idx_event_date;

-- +goose StatementEnd
