-- +goose Up
CREATE TABLE IF NOT EXISTS venues(
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    total_seats INTEGER NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
    deleted_at TIMESTAMPTZ,
  );
CREATE INDEX idx_venues_city ON venues(city);
-- +goose Down
DROP TABLE venues;
