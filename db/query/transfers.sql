-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account,
    to_account,
    amount
) VALUES ($1, $2, $3) RETURNING *;

-- name: GetTransfer :many
SELECT * FROM transfers
WHERE from_account = $1 AND to_account = $2 LIMIT $3;

-- name: DeleteTransfer :one
DELETE FROM transfers WHERE id = $1 RETURNING id;


