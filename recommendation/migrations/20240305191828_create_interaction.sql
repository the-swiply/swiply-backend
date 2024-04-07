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

CREATE INDEX IF NOT EXISTS idx_interaction_from ON interaction ("from");
CREATE INDEX IF NOT EXISTS idx_interaction_to_positive ON interaction ("to", positive);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS interaction;

DROP INDEX IF EXISTS idx_interaction_from;
DROP INDEX IF EXISTS idx_interaction_to_positive;
-- +goose StatementEnd
