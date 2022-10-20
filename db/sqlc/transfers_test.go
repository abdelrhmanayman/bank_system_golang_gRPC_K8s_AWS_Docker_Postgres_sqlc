package db

import (
	"banksystem/util"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createRandomTransferData(from_account, to_account int64) CreateTransferParams {
	return CreateTransferParams{
		FromAccount: from_account,
		ToAccount:   to_account,
		Amount:      util.CreateRandomBalance(),
	}
}

func createTwoTestAccounts(t *testing.T) (Account, Account) {
	assert := assert.New(t)

	accountOneData := CreateAccountTestArgs()
	accountTwoData := CreateAccountTestArgs()

	account1, err1 := CreateRandomAccount(accountOneData)
	account2, err2 := CreateRandomAccount(accountTwoData)

	assert.NotEqualValues(account1.ID, account2.ID)

	assert.NoError(err1)
	assert.NotEmpty(account1)

	assert.NoError(err2)
	assert.NotEmpty(account2)

	return account1, account2
}

func TestCreateTransfer(t *testing.T) {
	assert := assert.New(t)

	account1, account2 := createTwoTestAccounts(t)

	transfer, err := testQueries.CreateTransfer(context.Background(), createRandomTransferData(account1.ID, account2.ID))

	assert.NoError(err)
	assert.NotEmpty(transfer)
}

func TestGetTransfer(t *testing.T) {
	assert := assert.New(t)

	account1, account2 := createTwoTestAccounts(t)

	createdTransfer, createError := testQueries.CreateTransfer(context.Background(), createRandomTransferData(account1.ID, account2.ID))

	assert.NoError(createError)
	assert.NotEmpty(createdTransfer)

	transfers, err := testQueries.GetTransfer(context.Background(), GetTransferParams{FromAccount: account1.ID, ToAccount: account2.ID, Limit: 1})

	assert.NoError(err)
	assert.NotEmpty(transfers)
	assert.Len(transfers, 1)
}

func TestDeleteTransfer(t *testing.T) {
	assert := assert.New(t)

	account1, account2 := createTwoTestAccounts(t)

	createdTransfer, createError := testQueries.CreateTransfer(context.Background(), createRandomTransferData(account1.ID, account2.ID))

	assert.NoError(createError)
	assert.NotEmpty(createdTransfer)

	id, err := testQueries.DeleteTransfer(context.Background(), createdTransfer.ID)

	assert.NoError(err)
	assert.NotEmpty(id)
	assert.Equal(id, createdTransfer.ID)
}
