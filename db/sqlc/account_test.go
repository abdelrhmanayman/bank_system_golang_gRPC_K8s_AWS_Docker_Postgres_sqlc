package db

import (
	"banksystem/util"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateAccountTestArgs() CreateAccountParams {
	return CreateAccountParams{
		Owner:    util.CreateRandomOwner(),
		Balance:  util.CreateRandomBalance(),
		Currency: util.CreateRandomCurrency(),
	}
}

func CreateRandomAccount(accountArgs CreateAccountParams) (Account, error) {
	return testQueries.CreateAccount(context.Background(), accountArgs)
}

func TestCreateAccount(t *testing.T) {
	assert := assert.New(t)
	accountArgs := CreateAccountTestArgs()
	account, err := CreateRandomAccount(accountArgs)

	assert.NoError(err)
	assert.NotEmpty(account)
	assert.Equal(accountArgs.Balance, account.Balance)
	assert.Equal(accountArgs.Owner, account.Owner)
	assert.Equal(accountArgs.Currency, account.Currency)
}

func TestGetAccount(t *testing.T) {
	assert := assert.New(t)
	accountArgs := CreateAccountTestArgs()
	randomAccount, creationError := CreateRandomAccount(accountArgs)

	assert.NoError(creationError)

	account, err := testQueries.GetAccount(context.Background(), randomAccount.ID)

	assert.NoError(err)
	assert.NotEmpty(account)
	assert.Equal(randomAccount.ID, account.ID)
	assert.Equal(accountArgs.Balance, account.Balance)
	assert.Equal(accountArgs.Owner, account.Owner)
	assert.Equal(accountArgs.Currency, account.Currency)

}

func TestListAccounts(t *testing.T) {
	listArgs := ListAccountsParams{
		Limit:  2,
		Offset: 0,
	}
	assert := assert.New(t)
	accounts, err := testQueries.ListAccounts(context.Background(), listArgs)

	assert.NoError(err)
	assert.NotEmpty(accounts)
	assert.Len(accounts, 2)
}

func TestDeleteAccount(t *testing.T) {
	assert := assert.New(t)
	accountArgs := CreateAccountTestArgs()
	randomAccount, creationError := CreateRandomAccount(accountArgs)

	assert.NoError(creationError)

	err := testQueries.DeleteAccount(context.Background(), randomAccount.ID)

	assert.NoError(err)
}

func TestUpdateAccount(t *testing.T) {
	assert := assert.New(t)
	accountArgs := CreateAccountTestArgs()
	randomAccount, creationError := CreateRandomAccount(accountArgs)

	assert.NoError(creationError)

	updatedBalanceValue := util.CreateRandomBalance()

	updateArgs := UpdateAccountParams{
		ID:      randomAccount.ID,
		Balance: updatedBalanceValue,
	}

	updateAccount, err := testQueries.UpdateAccount(context.Background(), updateArgs)

	assert.NoError(err)
	assert.NotEmpty(updateAccount)
	assert.Equal(updateAccount.Balance, updatedBalanceValue)
}
