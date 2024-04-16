-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event
(
    id          bigserial PRIMARY KEY,
    owner       uuid,
    title       varchar(64),
    description text,
    chat_id     bigserial,
    date        timestamp
);

CREATE TABLE IF NOT EXISTS event_user_status
(
    id       bigserial PRIMARY KEY,
    user_id  uuid,
    event_id bigint,
    status   varchar(32)
);

ALTER TABLE event_user_status ADD UNIQUE (user_id, event_id);

CREATE INDEX IF NOT EXISTS idx_event_owner ON event (owner);
CREATE INDEX IF NOT EXISTS idx_event_title ON event (title);
CREATE INDEX IF NOT EXISTS idx_event_date ON event (date);

CREATE INDEX IF NOT EXISTS idx_event_user_status_user_id ON event_user_status (user_id);
CREATE INDEX IF NOT EXISTS idx_event_user_status_event_id ON event_user_status (event_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;

DROP INDEX IF EXISTS idx_event_owner;
DROP INDEX IF EXISTS idx_event_title;
DROP INDEX IF EXISTS idx_event_date;

DROP INDEX IF EXISTS idx_event_user_status_user_id;
DROP INDEX IF EXISTS idx_event_user_status_event_id;

-- +goose StatementEnd
