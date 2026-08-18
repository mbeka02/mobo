-- name: CreateEmailOutbox :one
INSERT INTO email_outbox (
    id, event_type, schema_version, aggregate_type, aggregate_id,
    recipient, template, payload
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ClaimEmailOutbox :many
WITH candidates AS (
    SELECT id
    FROM email_outbox
    WHERE status = 'pending'
      AND (lease_expires_at IS NULL OR lease_expires_at < now())
    ORDER BY created_at
    FOR UPDATE SKIP LOCKED
    LIMIT $1
)
UPDATE email_outbox AS outbox
SET status = 'processing',
    lease_owner = $2,
    lease_expires_at = $3,
    attempt_count = attempt_count + 1,
    updated_at = now()
FROM candidates
WHERE outbox.id = candidates.id
RETURNING outbox.*;

-- name: MarkEmailOutboxPublished :execrows
UPDATE email_outbox
SET status = 'published',
    published_at = now(),
    lease_owner = NULL,
    lease_expires_at = NULL,
    last_error = NULL,
    updated_at = now()
WHERE id = $1 AND lease_owner = $2 AND status = 'processing';

-- name: ReleaseEmailOutbox :execrows
UPDATE email_outbox
SET status = 'pending',
    lease_owner = NULL,
    lease_expires_at = NULL,
    last_error = $3,
    updated_at = now()
WHERE id = $1 AND lease_owner = $2 AND status = 'processing';
