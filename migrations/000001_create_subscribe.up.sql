CREATE TABLE subscribe (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    telegram_id BIGINT NOT NULL,
    city VARCHAR(255) NOT NULL,
    filters JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_subscribe_telegram_id ON subscribe (telegram_id);
CREATE INDEX idx_subscribe_city ON subscribe (city);
ALTER TABLE subscribe ADD CONSTRAINT unique_telegram_city UNIQUE (telegram_id, city);