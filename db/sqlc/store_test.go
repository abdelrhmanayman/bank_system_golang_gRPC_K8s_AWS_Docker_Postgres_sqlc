package db

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransferTx(t *testing.T) {
	assert := assert.New(t)

	store := NewStore(DBTestConnection)

	transactionResults := make(chan TransferTxResults)
	transactionErrors := make(chan error)

	numberOfTestTransactions := 10
	testAmount := int64(200)

	user1 := createUserTest(t)
	account1Args := CreateAccountTestArgs(user1.Username)

	user2 := createUserTest(t)
	account2Args := CreateAccountTestArgs(user2.Username)

	account1, _ := CreateRandomAccount(account1Args)
	account2, _ := CreateRandomAccount(account2Args)

	for i := 0; i < numberOfTestTransactions; i++ {
		go func() {
			result, err := store.ExecTransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        testAmount,
			})

			transactionResults <- result
			transactionErrors <- err

		}()
	}

	for i := 0; i < numberOfTestTransactions; i++ {
		transactionTestResult := <-transactionResults
		transactionTestError := <-transactionErrors

		assert.NoError(transactionTestError)
		assert.NotEmpty(transactionTestResult)

		// test transaction transfer
		resultTransfer := transactionTestResult.Transfer

		assert.NotEmpty(resultTransfer)
		assert.NotEmpty(resultTransfer.CreatedAt)
		assert.NotZero(resultTransfer.ID)
		assert.Equal(resultTransfer.Amount, testAmount)
		assert.Equal(resultTransfer.FromAccount, account1.ID)
		assert.Equal(resultTransfer.ToAccount, account2.ID)

		// test transaction entries

		fromEntry := transactionTestResult.FromEntry

		assert.NotEmpty(fromEntry)
		assert.NotZero(fromEntry.ID)
		assert.NotZero(fromEntry.CreatedAt)
		assert.Equal(fromEntry.Amount, -testAmount)
		assert.Equal(fromEntry.AccountID, transactionTestResult.FromAccount.ID)

		toEntry := transactionTestResult.ToEntry

		assert.NotEmpty(toEntry)
		assert.NotZero(toEntry.ID)
		assert.NotZero(toEntry.CreatedAt)
		assert.Equal(toEntry.Amount, testAmount)
		assert.Equal(toEntry.AccountID, transactionTestResult.ToAccount.ID)

		// test accounts after transaction
		assert.NotEmpty(transactionTestResult.FromAccount)
		assert.NotEmpty(transactionTestResult.ToAccount)

	}

}

func TestCreateTransferDeadlockTx(t *testing.T) {
	assert := assert.New(t)

	store := NewStore(DBTestConnection)

	transactionResults := make(chan TransferTxResults)
	transactionErrors := make(chan error)

	numberOfTestTransactions := 10
	testAmount := int64(200)

	user1 := createUserTest(t)
	account1Args := CreateAccountTestArgs(user1.Username)

	user2 := createUserTest(t)
	account2Args := CreateAccountTestArgs(user2.Username)

	account1, _ := CreateRandomAccount(account1Args)
	account2, _ := CreateRandomAccount(account2Args)

	fromAccounts := []int64{account1.ID, account2.ID}
	toAccounts := []int64{account2.ID, account1.ID}

	for i := 0; i < numberOfTestTransactions; i++ {
		counter := int64(i)
		go func() {
			counter = atomic.LoadInt64(&counter)
			result, err := store.ExecTransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccounts[(counter % 2)],
				ToAccountId:   toAccounts[(counter % 2)],
				Amount:        testAmount,
			})

			transactionResults <- result
			transactionErrors <- err

		}()
	}

	for i := 0; i < numberOfTestTransactions; i++ {
		transactionTestResult := <-transactionResults
		transactionTestError := <-transactionErrors

		assert.NoError(transactionTestError)
		assert.NotEmpty(transactionTestResult)

		// test transaction transfer
		resultTransfer := transactionTestResult.Transfer

		assert.NotEmpty(resultTransfer)
		assert.NotEmpty(resultTransfer.CreatedAt)
		assert.NotZero(resultTransfer.ID)
		assert.Equal(resultTransfer.Amount, testAmount)
		// Implement contains function for array in utils VI

		// test transaction entries

		fromEntry := transactionTestResult.FromEntry

		assert.NotEmpty(fromEntry)
		assert.NotZero(fromEntry.ID)
		assert.NotZero(fromEntry.CreatedAt)
		assert.Equal(fromEntry.Amount, -testAmount)
		assert.Equal(fromEntry.AccountID, transactionTestResult.FromAccount.ID)

		toEntry := transactionTestResult.ToEntry

		assert.NotEmpty(toEntry)
		assert.NotZero(toEntry.ID)
		assert.NotZero(toEntry.CreatedAt)
		assert.Equal(toEntry.Amount, testAmount)
		assert.Equal(toEntry.AccountID, transactionTestResult.ToAccount.ID)

		// test accounts after transaction
		assert.NotEmpty(transactionTestResult.FromAccount)
		assert.NotEmpty(transactionTestResult.ToAccount)

	}

}
