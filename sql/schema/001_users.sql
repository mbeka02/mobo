-- +goose Up 
CREATE TABLE IF NOT EXISTS users(
id UUID PRIMARY KEY,
email varchar(256) UNIQUE NOT NULL,
telephone_number varchar(16) NOT NULL,
password_hash varchar NOT NULL,
full_name varchar(256) NOT NULL,
profile_image_url varchar DEFAULT '',
user_name varchar(30),
created_at timestamptz NOT NULL DEFAULT (now()),
updated_at timestamptz DEFAULT (now()),
verified_at timestamptz,
deleted_at timestamptz
)
CREATE UNIQUE INDEX idx_users_email ON users(email)

-- +goose Down 
DROP TABLE users;
