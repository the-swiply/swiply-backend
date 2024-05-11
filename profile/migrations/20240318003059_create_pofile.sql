-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS interest
(
    id         bigserial PRIMARY KEY,
    definition text
);

CREATE TABLE IF NOT EXISTS organization
(
    id      bigserial PRIMARY KEY,
    "name"  text,
    pattern text
);

CREATE TYPE subscription_type AS ENUM ('STANDARD', 'PRIMARY');
CREATE TYPE gender_type AS ENUM ('MALE', 'FEMALE');

CREATE TABLE IF NOT EXISTS profile
(
    id            uuid               PRIMARY KEY,
    email         text               UNIQUE,
    "name"        text,
    city          text,
    "work"        text,
    education     text,
    is_blocked    bool,
    interests     bigint[],
    birth_day     timestamp,
    gender        gender_type,
    info          text,
    subscription  subscription_type,
    location_lat  double precision,
    location_long double precision,
    updated_at    timestamp
);

CREATE TABLE IF NOT EXISTS user_organization
(
    id              bigserial PRIMARY KEY,
    profile_id      uuid      REFERENCES profile (id),
    "name"          text,
    organization_id bigint    REFERENCES organization (id),
    email           text      UNIQUE,
    is_valid        bool
);

CREATE TYPE interaction_type AS ENUM ('LIKE', 'DISLIKE');

CREATE TABLE IF NOT EXISTS interaction
(
    id         bigserial            PRIMARY KEY,
    "from"     uuid                 REFERENCES profile (id),
    "to"       uuid                 REFERENCES profile (id),
    "type"     interaction_type,
    created_at timestamp
);

CREATE TABLE IF NOT EXISTS photo
(
    id         uuid   PRIMARY KEY,
    photo_ids  uuid[]
);

CREATE INDEX IF NOT EXISTS idx_interaction_from ON interaction ("from");
CREATE INDEX IF NOT EXISTS idx_interaction_to ON interaction ("to");
CREATE INDEX IF NOT EXISTS idx_interaction_created_at ON interaction (created_at);
CREATE INDEX IF NOT EXISTS idx_profile_updated_at ON profile (updated_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_interaction_to;
DROP INDEX IF EXISTS idx_interaction_from;
DROP INDEX IF EXISTS idx_interaction_created_at;
DROP INDEX IF EXISTS idx_profile_updated_at;

DROP TABLE IF EXISTS photo;
DROP TABLE IF EXISTS interaction;

DROP TYPE IF EXISTS interaction_type;

DROP TABLE IF EXISTS user_organization;

DROP TABLE IF EXISTS profile;

DROP TYPE IF EXISTS gender_type;
DROP TYPE IF EXISTS subscription_type;

DROP TABLE IF EXISTS organization;

DROP TABLE IF EXISTS interest;
-- +goose StatementEnd
