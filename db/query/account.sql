-- name: CreateAccount :one
INSERT INTO account (
    owner,
    balance,
    currency
) VALUES ($1, $2, $3) RETURNING *;


-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM account WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM account ORDER BY id LIMIT $1 OFFSET $2;

-- name: DeleteAccount :exec
DELETE FROM account WHERE id = $1;

-- name: UpdateAccount :one
UPDATE account SET balance = $2 WHERE id = $1 RETURNING *;

-- name: AddMoneyToAccount :one
UPDATE account SET balance = balance + $1 WHERE id = $2 RETURNING *;