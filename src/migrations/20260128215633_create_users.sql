-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                                     id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     username      TEXT UNIQUE NOT NULL,
                                     password_hash TEXT NOT NULL,
                                     created_at    timestamp WITH TIME ZONE

);

-- Опционально: индекс (если часто ищешь по username)
CREATE INDEX idx_users_username ON users(username);


-- +goose Down
DROP TABLE IF EXISTS users;