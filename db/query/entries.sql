-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES ($1, $2) RETURNING *;

-- name: GetAccountEntries :many
SELECT * FROM entries
WHERE account_id = $1 LIMIT $2;

-- name: DeleteEntry :one
DELETE FROM entries WHERE id = $1 RETURNING id;


