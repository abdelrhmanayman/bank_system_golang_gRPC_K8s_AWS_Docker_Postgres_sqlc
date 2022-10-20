// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: entries.sql

package db

import (
	"context"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES ($1, $2) RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteEntry = `-- name: DeleteEntry :one
DELETE FROM entries WHERE id = $1 RETURNING id
`

func (q *Queries) DeleteEntry(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, deleteEntry, id)
	err := row.Scan(&id)
	return id, err
}

const getEntry = `-- name: GetEntry :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $1 LIMIT $2
`

type GetEntryParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
}

func (q *Queries) GetEntry(ctx context.Context, arg GetEntryParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getEntry, arg.AccountID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
