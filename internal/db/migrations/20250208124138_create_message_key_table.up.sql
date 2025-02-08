CREATE TABLE message_key (
    id UUID PRIMARY KEY,
    message_id UUID REFERENCES message(id),
    part_index INT NOT NULL,
    secret_part CHAR(64),
    updated_at TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE message_key is 'таблица для хранения частей ключа сообщения с указанием порядка частей';

COMMENT ON COLUMN message_key.id is 'уникальный идентификатор части ключа';

COMMENT ON COLUMN message_key.message_id is 'связь с сообщением';

COMMENT ON COLUMN message_key.part_index is 'порядковый номер части ключа';

COMMENT ON COLUMN message_key.secret_part is 'часть ключа, введенная пользователем, в формате HEX';

COMMENT ON COLUMN message_key.updated_at is 'дата последнего ввода части ключа пользователем';
