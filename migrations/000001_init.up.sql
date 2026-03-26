CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    master_hash BYTEA,
    master_salt BYTEA
);

CREATE TABLE IF NOT EXISTS passwords (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    service VARCHAR(64) NOT NULL,
    encrypted_data BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    master_check TEXT DEFAULT 'OK',
    UNIQUE(user_id, service)
);