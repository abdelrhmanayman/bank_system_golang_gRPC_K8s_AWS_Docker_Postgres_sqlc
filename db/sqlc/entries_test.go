package db

import (
	"banksystem/util"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEntry(t *testing.T) {
	assert := assert.New(t)

	user := createUserTest(t)
	accountOneData := CreateAccountTestArgs(user.Username)
	account1, err1 := CreateRandomAccount(accountOneData)

	assert.NoError(err1)
	assert.NotEmpty(account1)

	entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{AccountID: account1.ID, Amount: util.CreateRandomBalance()})

	assert.NoError(err)
	assert.NotEmpty(entry)
}

func TestGetEntry(t *testing.T) {
	assert := assert.New(t)

	user := createUserTest(t)
	accountOneData := CreateAccountTestArgs(user.Username)
	account1, err1 := CreateRandomAccount(accountOneData)

	assert.NoError(err1)
	assert.NotEmpty(account1)

	createdEntry, createErr := testQueries.CreateEntry(context.Background(), CreateEntryParams{AccountID: account1.ID, Amount: util.CreateRandomBalance()})

	assert.NoError(createErr)
	assert.NotEmpty(createdEntry)

	entries, err := testQueries.GetAccountEntries(context.Background(), GetAccountEntriesParams{
		AccountID: account1.ID,
		Limit:     1,
	})

	assert.NoError(err)
	assert.NotEmpty(entries)
	assert.Len(entries, 1)

}

func TestDeleteEntry(t *testing.T) {
	assert := assert.New(t)

	user := createUserTest(t)
	accountOneData := CreateAccountTestArgs(user.Username)
	account1, err1 := CreateRandomAccount(accountOneData)

	assert.NoError(err1)
	assert.NotEmpty(account1)

	createdEntry, createErr := testQueries.CreateEntry(context.Background(), CreateEntryParams{AccountID: account1.ID, Amount: util.CreateRandomBalance()})

	assert.NoError(createErr)
	assert.NotEmpty(createdEntry)

	id, err := testQueries.DeleteEntry(context.Background(), createdEntry.ID)

	assert.NoError(err)
	assert.NotEmpty(id)
	assert.Equal(id, createdEntry.ID)
}
