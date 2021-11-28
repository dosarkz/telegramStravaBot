CREATE TABLE users
(
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(255) NOT NULL,
    telegram_id     integer NOT NULL,
    created_at      TIMESTAMP           NOT NULL,
    updated_at      TIMESTAMP           NOT NULL,
    deleted_at      TIMESTAMP NULL
);

CREATE UNIQUE INDEX idx_users_telegram_id ON users (telegram_id);