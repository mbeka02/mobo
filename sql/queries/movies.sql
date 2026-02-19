-- name: AddMovie :one
INSERT INTO movies (
title,description,runtime,genre,age_rating,director,poster_url,release_date 
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- name: GetMovieById :one 
SELECT id,title,description,runtime,genre,age_rating,director,poster_url,release_date
,created_at,updated_at
FROM movies 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetMovies :many
SELECT id,title,description,runtime,genre,age_rating,director,poster_url,release_date
,created_at,updated_at
FROM movies 
WHERE  deleted_at IS NULL
ORDER BY release_date DESC 
LIMIT $1 OFFSET $2;

-- name: DeleteMovie :exec
UPDATE movies SET deleted_at=now() WHERE id=$1;

