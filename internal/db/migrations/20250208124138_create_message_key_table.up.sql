CREATE TABLE message_key (
    id UUID PRIMARY KEY,
    message_id UUID REFERENCES message(id),
    secret_part BYTEA,
    updated_at TIMESTAMPTZ NOT NULL,
    CHECK (octet_length(secret_part) <= 256)
);

COMMENT ON TABLE message_key is 'таблица для хранения частей ключа сообщения с указанием порядка частей';

COMMENT ON COLUMN message_key.id is 'уникальный идентификатор части ключа';

COMMENT ON COLUMN message_key.message_id is 'связь с сообщением';

COMMENT ON COLUMN message_key.secret_part is 'часть ключа, введенная пользователем';

COMMENT ON COLUMN message_key.updated_at is 'дата последнего ввода части ключа пользователем';
