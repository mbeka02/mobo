-- +goose Up 
CREATE TABLE IF NOT EXISTS showtimes(
id BIGSERIAL PRIMARY KEY,
movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
start_time TIMESTAMPTZ NOT NULL,
end_time TIMESTAMPTZ NOT NULL,
available_seats INTEGER NOT NULL,
price_per_seat NUMERIC(10,2) NOT NULL,
venue_id INT NOT NULL REFERENCES venues(id)
created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
updated_at TIMESTAMPTZ DEFAULT (now()),
deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_showtimes_movie_id ON showtimes(movie_id);
ALTER TABLE showtimes ADD CONSTRAINT chk_showtime_times 
    CHECK (end_time > start_time);

ALTER TABLE showtimes ADD CONSTRAINT chk_available_seats 
    CHECK (available_seats >= 0);

-- +goose Down 
DROP TABLE showtimes;
