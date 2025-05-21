-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = $2, revoked_at = $3
WHERE token = $1;