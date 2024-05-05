-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS oracle_prediction
(
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        uuid,
    recommendation uuid,
    score          double precision
);

CREATE INDEX IF NOT EXISTS idx_oracle_prediction_user_id ON oracle_prediction (user_id);
CREATE INDEX IF NOT EXISTS idx_oracle_prediction_score ON oracle_prediction (score);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oracle_prediction;

DROP INDEX IF EXISTS idx_oracle_prediction_user_id;
DROP INDEX IF EXISTS idx_oracle_prediction_score;
-- +goose StatementEnd