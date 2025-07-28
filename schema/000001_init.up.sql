CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(200) NOT NULL UNIQUE,
    email VARCHAR(200) UNIQUE,
    telegram_id BIGINT UNIQUE,
    password_hash VARCHAR(256) NOT NULL
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    refresh_token TEXT NOT NULL,
    fingerprint VARCHAR(256) NOT NULL,
    UNIQUE (user_id, fingerprint)
);