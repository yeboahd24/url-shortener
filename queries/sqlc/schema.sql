-- schema.sql
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);

CREATE TABLE urls (
    short_id VARCHAR(10) PRIMARY KEY,
    long_url TEXT NOT NULL,
    user_id UUID REFERENCES users(user_id),
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    click_limit INTEGER
);

CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    short_id VARCHAR(10) REFERENCES urls(short_id),
    ip_address VARCHAR(45),
    user_agent TEXT,
    clicked_at TIMESTAMP NOT NULL
);

CREATE TABLE api_keys (
    key UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMP NOT NULL
);
