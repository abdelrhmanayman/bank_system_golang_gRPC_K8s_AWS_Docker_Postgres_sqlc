package db

import (
	"banksystem/util"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateRandomUser() CreateUserParams {
	hashedPassword, _ := util.HashPassword(util.GenerateRandomString(8))

	return CreateUserParams{
		Username:  util.CreateRandomOwner(),
		HashedPwd: hashedPassword,
		Email:     util.CreateRandomEmail(),
	}
}

func createUserTest(t *testing.T) User {
	assert := assert.New(t)

	userArgs := CreateRandomUser()

	user, err := testQueries.CreateUser(context.Background(), userArgs)

	assert.NoError(err)
	assert.NotEmpty(user)

	assert.Equal(userArgs.Username, user.Username)
	assert.Equal(userArgs.Email, user.Email)
	assert.Equal(userArgs.HashedPwd, user.HashedPwd)
	assert.NotEmpty(user.PasswordChangedAt)
	assert.NotEmpty(user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createUserTest(t)
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	createdUser := createUserTest(t)
	user, err := testQueries.GetUser(context.Background(), createdUser.Username)

	assert.NoError(err)
	assert.NotEmpty(user)
	assert.Equal(createdUser.Email, user.Email)
	assert.Equal(createdUser.HashedPwd, user.HashedPwd)
	assert.Equal(createdUser.Username, user.Username)
	assert.NotEmpty(user.CreatedAt)
	assert.NotEmpty(user.PasswordChangedAt)

}
