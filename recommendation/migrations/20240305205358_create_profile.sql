-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS profile
(
    id                 uuid PRIMARY KEY,
    interests          bigint[],
    birthday           timestamp,
    gender             varchar(16),
    info               text,
    subscription_type  varchar(32),
    location_lat       double precision,
    location_lon       double precision,
    updated_at         timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS profile;
-- +goose StatementEnd
