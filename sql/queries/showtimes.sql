-- name: CreateShowtime :one
INSERT INTO showtimes (movie_id, start_time, end_time, available_seats, price_per_seat, venue_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetShowtimeById :one
SELECT * FROM showtimes WHERE id = $1 AND deleted_at IS NULL;

-- name: GetShowtimesByMovie :many
SELECT s.*, v.name as venue_name, v.city as venue_city
FROM showtimes s
JOIN venues v ON v.id = s.venue_id
WHERE s.movie_id = $1
  AND s.start_time > now()
  AND s.deleted_at IS NULL
ORDER BY s.start_time ASC;

-- name: GetShowtimesAdmin :many
SELECT s.*, m.title as movie_title, v.name as venue_name
FROM showtimes s
JOIN movies m ON m.id = s.movie_id
JOIN venues v ON v.id = s.venue_id
WHERE s.deleted_at IS NULL
ORDER BY s.start_time DESC
LIMIT $1 OFFSET $2;

-- name: UpdateShowtime :one
UPDATE showtimes SET
  start_time = COALESCE(sqlc.narg('start_time'), start_time),
  end_time = COALESCE(sqlc.narg('end_time'), end_time),
  available_seats = COALESCE(sqlc.narg('available_seats'), available_seats),
  price_per_seat = COALESCE(sqlc.narg('price_per_seat'), price_per_seat),
  venue_id = COALESCE(sqlc.narg('venue_id'), venue_id),
  updated_at = now()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: DeleteShowtime :exec
UPDATE showtimes SET deleted_at = now() WHERE id = $1;
