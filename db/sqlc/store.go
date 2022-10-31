package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	// return error if transaction can't begin
	if err != nil {
		return err
	}

	query := New(tx)
	err = fn(query)

	if err != nil {
		// rollback the transaction
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return fmt.Errorf("transaction error %v, rollback error %v", err, rollbackError)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResults struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	ToEntry     Entry    `json:"to_entry"`
	FromEntry   Entry    `json:"from_entry"`
}

func addMoney(
	q *Queries,
	ctx context.Context,
	accOneId int64,
	accOneAmount int64,
	accTwoId int64,
	accTwoAmount int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddMoneyToAccount(ctx, AddMoneyToAccountParams{
		ID:      accOneId,
		Balance: accOneAmount,
	})

	if err != nil {
		return
	}

	account2, err = q.AddMoneyToAccount(ctx, AddMoneyToAccountParams{
		ID:      accTwoId,
		Balance: accTwoAmount,
	})

	if err != nil {
		return
	}

	return
}

func (store *Store) ExecTransferTx(ctx context.Context, args TransferTxParams) (TransferTxResults, error) {
	var result TransferTxResults

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// fromAccount, fromAccountGetError := q.GetAccountForUpdate(ctx, args.FromAccountId)

		// if fromAccountGetError != nil {
		// 	return fromAccountGetError
		// }

		// toAccount, toAccountGetError := q.GetAccountForUpdate(ctx, args.ToAccountId)

		// if toAccountGetError != nil {
		// 	return toAccountGetError
		// }

		// if balanceDiff := (fromAccount.Balance - args.Amount); balanceDiff < 0 {
		// 	return errors.New("Insufficient balance to make the transfer")
		// }

		if args.FromAccountId < args.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(q, ctx, args.FromAccountId, -args.Amount, args.ToAccountId, args.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(q, ctx, args.ToAccountId, args.Amount, args.FromAccountId, -args.Amount)
		}

		if result.FromAccount.Balance < 0 {
			return errors.New("Insufficient balance to make the transfer")
		}

		if err != nil {
			return err
		}

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccount: args.FromAccountId,
			ToAccount:   args.ToAccountId,
			Amount:      args.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountId,
			Amount:    -args.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountId,
			Amount:    args.Amount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
