-- name: AddMovie :one
INSERT INTO movies (
title,description,runtime,genre,age_rating,director,poster_url,release_date 
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- admin listing , gets all the movies even if they aren't showing.
-- name: GetMoviesAdmin :many
SELECT id,title,description,runtime,genre,age_rating,director,poster_url,release_date
,created_at,updated_at
FROM movies 
WHERE  deleted_at IS NULL
ORDER BY release_date DESC 
LIMIT $1 OFFSET $2;

-- name: DeleteMovie :exec
UPDATE movies SET deleted_at=now() WHERE id=$1;

-- public listing: only movies with at least one future showtime
-- name: GetMoviesPublic :many
SELECT DISTINCT m.* FROM movies m
JOIN showtimes s ON s.movie_id = m.id
WHERE m.deleted_at IS NULL
  AND s.start_time > now()
  AND s.deleted_at IS NULL
ORDER BY s.start_time ASC
LIMIT $1 OFFSET $2;

-- name: GetMovieById :one 
SELECT id,title,description,runtime,genre,age_rating,director,poster_url,release_date
,created_at,updated_at
FROM movies 
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateMovie :one
UPDATE movies SET
  title = COALESCE(sqlc.narg('title'), title),
  description = COALESCE(sqlc.narg('description'), description),
  runtime = COALESCE(sqlc.narg('runtime'), runtime),
  genre = COALESCE(sqlc.narg('genre'), genre),
  age_rating = COALESCE(sqlc.narg('age_rating'), age_rating),
  director = COALESCE(sqlc.narg('director'), director),
  poster_url = COALESCE(sqlc.narg('poster_url'), poster_url),
  release_date = COALESCE(sqlc.narg('release_date'), release_date),
  updated_at = now()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;
