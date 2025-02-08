CREATE TABLE message (
    id UUID PRIMARY KEY,
    creator_name VARCHAR(16) NOT NULL,
    encrypted_content VARCHAR(1024) NOT NULL,
    key_hash CHAR(64) NOT NULL
);

COMMENT ON TABLE message is 'таблица для хранения зашифрованных последних сообщений пользователей EverLock';

COMMENT ON COLUMN message.id IS 'уникальный идентификатор сообщения';

COMMENT ON COLUMN message.creator_name IS 'ник создателя сообщения';

COMMENT ON COLUMN message.encrypted_content IS 'зашифрованное сообщение';

COMMENT ON COLUMN message.key_hash IS 'хеш оригинального ключа (SHA-256)';
