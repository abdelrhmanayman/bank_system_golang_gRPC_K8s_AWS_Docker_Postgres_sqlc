// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"
)

type Querier interface {
	AddMoneyToAccount(ctx context.Context, arg AddMoneyToAccountParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	DeleteAccount(ctx context.Context, id int64) error
	DeleteEntry(ctx context.Context, id int64) (int64, error)
	DeleteTransfer(ctx context.Context, id int64) (int64, error)
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountEntries(ctx context.Context, arg GetAccountEntriesParams) ([]Entry, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetTransfer(ctx context.Context, arg GetTransferParams) ([]Transfer, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
