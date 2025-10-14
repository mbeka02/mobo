-- +goose Up 
CREATE TABLE reservations(
    id BIGSERIAL PRIMARY KEY,
    showtime_id BIGINT NOT NULL REFERENCES showtimes(id),
    user_id UUID NOT NULL REFERENCES users(id),
    number_of_seats INTEGER NOT NULL CHECK (number_of_seats > 0),
    total_cost NUMERIC(10,2) NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'pending',  -- pending, confirmed, expired, cancelled
    reserved_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    expires_at TIMESTAMPTZ, --reservation expires in 10 minutes if not paid for
    confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now())
    deleted_at TIMESTAMPTZ,
  );
CREATE INDEX idx_reservations_user_id ON reservations(user_id);
CREATE INDEX idx_reservations_showtime_id ON reservations(showtime_id);
-- +goose Down 
DROP TABLE reservations;
