-- name: CreateUser :one
INSERT INTO users(
    email, 
    full_name, 
    profile_image_url,
    verified_at
) VALUES($1, $2, $3, $4) 
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users 
SET password_hash = $2, 
    updated_at = now() 
WHERE id = $1 
RETURNING *;

-- name: GetUserByEmail :one 
SELECT * FROM users 
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserById :one
SELECT * FROM users 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetIdentityByProvider :one
SELECT * FROM user_identities 
WHERE provider = $1 
  AND provider_user_id = $2;

-- name: LinkIdentityToUser :one
INSERT INTO user_identities (
    user_id,
    provider,
    provider_user_id
) VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateLocalUser :one
INSERT INTO users(
    email, 
    telephone_number, 
    password_hash, 
    full_name,
    verified_at
) VALUES($1, $2, $3, $4, NULL) 
RETURNING *;

-- name: GetUserByProvider :one
SELECT u.* 
FROM users u
JOIN user_identities ui ON u.id = ui.user_id
WHERE ui.provider = $1 
  AND ui.provider_user_id = $2 
  AND u.deleted_at IS NULL;
