
-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, expires_at, user_id) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetRefreshTokenForUser :one
SELECT refresh_tokens.token, refresh_tokens.user_id FROM refresh_tokens
WHERE refresh_tokens.user_id = $1 AND 
    (refresh_tokens.revoked_at IS NOT NULL AND refresh_tokens.expires_at >= CURRENT_TIMESTAMP);

-- name: GetRefreshToken :one
SELECT refresh_tokens.token, refresh_tokens.user_id FROM refresh_tokens
WHERE refresh_tokens.token = $1 AND 
    (refresh_tokens.revoked_at IS NOT NULL AND refresh_tokens.expires_at >= CURRENT_TIMESTAMP);

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $2
WHERE token = $1;