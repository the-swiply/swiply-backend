-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS update_info
(
    id serial,
    entity varchar(64),
    last_update timestamp
);

INSERT INTO update_info (entity, last_update)
VALUES ('interaction', now()),
       ('profile', now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS update_info;
-- +goose StatementEnd
