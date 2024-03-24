-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS recommendation_history
(
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        uuid,
    recommendation uuid,
    dttm           timestamp
);

CREATE INDEX IF NOT EXISTS idx_recommendation_history_user_id ON recommendation_history (user_id);
CREATE INDEX IF NOT EXISTS idx_recommendation_history_dttm ON recommendation_history (dttm);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS recommendation_history;

DROP INDEX IF EXISTS idx_recommendation_history_user_id;
DROP INDEX IF EXISTS idx_recommendation_history_dttm;
-- +goose StatementEnd
