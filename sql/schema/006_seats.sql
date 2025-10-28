-- +goose Up
CREATE TABLE seats(
    id BIGSERIAL PRIMARY KEY,
    showtime_id BIGINT NOT NULL REFERENCES showtimes(id) ON DELETE RESTRICT,
    row_letter VARCHAR(1) NOT NULL,
    seat_number VARCHAR NOT NULL,
    reservation_id BIGINT REFERENCES reservations(id) ON DELETE SET NULL,
    reserved_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    UNIQUE(showtime_id, row_letter, seat_number)
);


-- +goose Down 
DROP TABLE seats;
