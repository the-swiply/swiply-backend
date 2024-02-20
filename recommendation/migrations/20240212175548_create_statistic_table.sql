-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS statistic
(
    id         uuid,
    user_id    uuid,
    like_ratio double precision,
    updated_at timestamp
);

CREATE INDEX IF NOT EXISTS idx_statistic_user_id ON statistic (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS statistic;

DROP INDEX IF EXISTS idx_statistic_user_id;
-- +goose StatementEnd
