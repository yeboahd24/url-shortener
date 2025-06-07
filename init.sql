-- Database initialization script for URL Shortener
-- This script creates all necessary tables and indexes

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- Create urls table
CREATE TABLE IF NOT EXISTS urls (
    short_id VARCHAR(10) PRIMARY KEY,
    long_url TEXT NOT NULL,
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    click_limit INTEGER,
    CONSTRAINT valid_click_limit CHECK (click_limit IS NULL OR click_limit > 0)
);

-- Create clicks table
CREATE TABLE IF NOT EXISTS clicks (
    id SERIAL PRIMARY KEY,
    short_id VARCHAR(10) REFERENCES urls(short_id) ON DELETE CASCADE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    clicked_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create api_keys table
CREATE TABLE IF NOT EXISTS api_keys (
    key UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id);
CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at);
CREATE INDEX IF NOT EXISTS idx_urls_expires_at ON urls(expires_at) WHERE expires_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_clicks_short_id ON clicks(short_id);
CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);
CREATE INDEX IF NOT EXISTS idx_clicks_ip_address ON clicks(ip_address);

CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_created_at ON api_keys(created_at);

-- Create a function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for users table
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data (optional - remove in production)
-- INSERT INTO users (user_id, username, email, created_at) 
-- VALUES ('550e8400-e29b-41d4-a716-446655440000', 'admin', 'admin@example.com', NOW())
-- ON CONFLICT (user_id) DO NOTHING;

-- INSERT INTO api_keys (key, user_id, created_at)
-- VALUES ('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', NOW())
-- ON CONFLICT (key) DO NOTHING;
