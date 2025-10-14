-- +goose Up 
CREATE TABLE IF NOT EXISTS movies(
id BIGSERIAL PRIMARY KEY,
title VARCHAR NOT NULL,
description VARCHAR NOT NULL,
runtime INTEGER NOT NULL, -- SHOULD BE IN SECONDS 
genre VARCHAR NOT NULL,
age_rating VARCHAR NOT NULL,
director VARCHAR NOT NULL,
poster_url VARCHAR NOT NULL,
release_date DATE NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
updated_at TIMESTAMPTZ DEFAULT (now()),
deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_movies_genre ON movies(genre);
CREATE INDEX idx_movies_director ON movies(director);

-- +goose Down
DROP TABLE MOVIES;
