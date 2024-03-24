-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS profile
(
    id         uuid PRIMARY KEY,
   -- TODO
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS profile;
-- +goose StatementEnd
