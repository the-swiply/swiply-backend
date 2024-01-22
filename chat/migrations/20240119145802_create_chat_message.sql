-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat
(
    id      bigserial PRIMARY KEY,
    members uuid[]
);

CREATE TABLE IF NOT EXISTS message
(
    id         uuid,
    chat_id    bigint REFERENCES chat (id),
    id_in_chat bigint,
    send_time  timestamp,
    content    text

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat;
DROP TABLE IF EXISTS message;
-- +goose StatementEnd
