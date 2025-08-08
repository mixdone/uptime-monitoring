CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(200) NOT NULL UNIQUE,
    email VARCHAR(200) UNIQUE,
    telegram_id BIGINT UNIQUE,
    password_hash VARCHAR(256) NOT NULL
);

CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at DATE NOT NULL,
    fingerprint VARCHAR(256) NOT NULL,
    UNIQUE (user_id, fingerprint)
);


CREATE TABLE monitors (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    type VARCHAR(10) NOT NULL DEFAULT 'http',  
    target TEXT NOT NULL,                        
    timeout INT NOT NULL DEFAULT 10 CHECK (timeout BETWEEN 1 AND 300),
    interval INT NOT NULL DEFAULT 60 CHECK (interval BETWEEN 10 AND 3600),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_checked_at TIMESTAMPTZ,
); 


CREATE TABLE monitor_specs (
    id BIGSERIAL PRIMARY KEY,
    monitor_id BIGINT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,

    request JSONB NOT NULL,
    expected_response JSONB
);


