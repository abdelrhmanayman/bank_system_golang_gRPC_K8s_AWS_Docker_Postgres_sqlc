package token

import (
	"banksystem/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasetoMaker(t *testing.T) {
	assert := assert.New(t)
	username := util.CreateRandomOwner()
	symmetricKey := util.GenerateRandomString(32)
	pMaker, err := CreateNewPasetoMaker([]byte(symmetricKey))

	assert.NoError(err)

	token, payload, err := pMaker.CreateToken(username, time.Minute)

	assert.NoError(err)
	assert.NotEmpty(payload)
	assert.NotEmpty(token)
	assert.Equal(payload.Username, username)

	payload, err = pMaker.VerifyToken(token)

	assert.NoError(err)
	assert.NotEmpty(payload)
	assert.Equal(payload.Username, username)

}

func TestExpiredPasetoToken(t *testing.T) {
	assert := assert.New(t)
	username := util.CreateRandomOwner()
	symmetricKey := util.GenerateRandomString(32)
	pMaker, err := CreateNewPasetoMaker([]byte(symmetricKey))

	assert.NoError(err)

	token, payload, err := pMaker.CreateToken(username, -time.Minute)

	assert.NoError(err)
	assert.NotEmpty(payload)
	assert.NotEmpty(token)
	assert.Equal(payload.Username, username)

	payload, err = pMaker.VerifyToken(token)

	assert.Errorf(err, ErrExpiredToken.Error())
	assert.Empty(payload)

}
