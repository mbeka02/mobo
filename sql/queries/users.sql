-- name: CreateLocalUser :one
INSERT INTO users(
    email, 
    telephone_number, 
    password_hash, 
    full_name, 
    auth_provider
) VALUES($1, $2, $3, $4, 'local') 
RETURNING *;

-- name: CreateOAuthUser :one
INSERT INTO users(
    email, 
    full_name, 
    auth_provider, 
    provider_user_id,
    profile_image_url,
    verified_at  -- OAuth users are pre-verified by the provider
) VALUES($1, $2, $3, $4, $5, now()) 
RETURNING *;

-- name: GetUserByEmail :one 
SELECT 
    id, email, telephone_number, full_name, 
    profile_image_url, user_name, auth_provider,
    created_at, updated_at, verified_at 
FROM users 
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByProvider :one
SELECT 
    id, email, telephone_number, full_name, 
    profile_image_url, user_name, auth_provider,
    created_at, updated_at, verified_at 
FROM users 
WHERE auth_provider = $1 
  AND provider_user_id = $2 
  AND deleted_at IS NULL;

