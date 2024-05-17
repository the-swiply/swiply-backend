-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification
(
    id           uuid PRIMARY KEY,
    device_token text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF NOT EXISTS notification;
-- +goose StatementEnd
