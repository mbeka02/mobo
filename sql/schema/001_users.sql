-- +goose Up 
CREATE EXTENSION IF NOT EXISTS pg_uuidv7;
CREATE TABLE IF NOT EXISTS users(
id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
email varchar(256) UNIQUE NOT NULL,
telephone_number varchar(16) NOT NULL,
password_hash varchar,
auth_provider varchar(50) DEFAULT 'local',
provider_user_id varchar(256),  
full_name varchar(256) NOT NULL,
profile_image_url varchar DEFAULT '',
user_name varchar(30) DEFAULT '',
created_at timestamptz NOT NULL DEFAULT (now()),
updated_at timestamptz DEFAULT (now()),
verified_at timestamptz,
deleted_at timestamptz,
-- Ensure password exists for local auth
CONSTRAINT password_required_for_local 
CHECK (auth_provider != 'local' OR password_hash IS NOT NULL),
CONSTRAINT unique_provider_user UNIQUE (auth_provider, provider_user_id)
);
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_provider ON users(auth_provider, provider_user_id);
-- +goose Down 
DROP TABLE users;
