-- +goose Up
-- +goose StatementBegin
CREATE TYPE meeting_status AS ENUM ('AWAITING_SCHEDULE', 'SCHEDULING', 'SCHEDULED');

CREATE TABLE IF NOT EXISTS meeting
(
    id              uuid           PRIMARY KEY,
    owner_id        uuid,
    member_id       uuid,
    "start"         timestamp,
    "end"           timestamp,
    organization_id bigint,
    status          meeting_status
);

CREATE INDEX IF NOT EXISTS idx_meeting_owner_id ON meeting (owner_id);
CREATE INDEX IF NOT EXISTS idx_meeting_start ON meeting ("start");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_meeting_owner_id;
DROP INDEX IF EXISTS idx_meeting_start;

DROP TABLE IF EXISTS meeting;

DROP TYPE IF EXISTS meeting_status;
-- +goose StatementEnd