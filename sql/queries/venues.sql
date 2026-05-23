-- name: CreateVenue :one
INSERT INTO venues (name, address, city, total_seats)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetVenues :many
SELECT * FROM venues WHERE deleted_at IS NULL ORDER BY name;

-- name: GetVenueById :one
SELECT * FROM venues WHERE id = $1 AND deleted_at IS NULL;
