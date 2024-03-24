-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS statistic
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid,
    like_ratio double precision,
    updated_at timestamp
);

CREATE INDEX IF NOT EXISTS idx_statistic_user_id ON statistic (user_id);
CREATE INDEX IF NOT EXISTS idx_statistic_like_ratio ON statistic (like_ratio);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS statistic;

DROP INDEX IF EXISTS idx_statistic_user_id;
DROP INDEX IF EXISTS idx_statistic_like_ratio;
-- +goose StatementEnd
