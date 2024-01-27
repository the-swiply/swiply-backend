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
    "from"     uuid,
    send_time  timestamp,
    content    text
);

CREATE INDEX IF NOT EXISTS idx_message_chat_id_id_in_chat ON message (chat_id, id_in_chat);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat;
DROP TABLE IF EXISTS message;

DROP INDEX IF EXISTS idx_message_chat_id_id_in_chat;
-- +goose StatementEnd
