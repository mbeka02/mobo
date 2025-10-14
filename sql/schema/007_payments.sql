-- +goose Up

CREATE TABLE IF NOT EXISTS payments(
    id BIGSERIAL PRIMARY KEY,
    reservation_id BIGINT NOT NULL REFERENCES reservations(id),
    amount NUMERIC(10,2) NOT NULL,
    payment_method VARCHAR NOT NULL,  -- card, cash, etc
    payment_status VARCHAR NOT NULL,  -- pending, completed, failed, refunded
    transaction_id VARCHAR, 
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    updated_at TIMESTAMPTZ DEFAULT (now())
    deleted_at TIMESTAMPTZ,
  );

-- +goose Down 
DROP TABLE payments;
