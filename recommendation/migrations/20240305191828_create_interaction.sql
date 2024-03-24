-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS interaction
(
    id         uuid PRIMARY KEY,
    "from"     uuid,
    "to"       uuid,
    positive   bool,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS interaction;
-- +goose StatementEnd
