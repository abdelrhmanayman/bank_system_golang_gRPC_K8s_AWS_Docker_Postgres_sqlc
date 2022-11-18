-- name: CreateSession :one
INSERT INTO sessions (
    username,
    refresh_token,
    user_agent,
    client_ip,
    expires_at,
    is_blocked,
    id
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: UpdateSession :one
UPDATE sessions SET is_blocked = $2, expires_at = $3 WHERE id = $1 RETURNING *;

