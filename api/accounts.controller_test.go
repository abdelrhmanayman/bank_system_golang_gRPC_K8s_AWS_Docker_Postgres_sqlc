package api

import (
	mockdb "banksystem/db/mock"
	db "banksystem/db/sqlc"
	"banksystem/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAccountController(t *testing.T) {
	assert := assert.New(t)
	account := createDummyAccount()

	ctrl := gomock.NewController(t)

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)

	server := SetupRoutes(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(err)

	server.router.ServeHTTP(recorder, request)

	assert.Equal(http.StatusOK, recorder.Code)
	matchAccountBodyCheck(t, account, recorder.Body)

}

func createDummyAccount() db.Account {
	return db.Account{
		ID:       util.CreateRandomBalance(),
		Balance:  util.CreateRandomBalance(),
		Owner:    util.CreateRandomOwner(),
		Currency: util.CreateRandomCurrency(),
	}
}

func matchAccountBodyCheck(t *testing.T, account db.Account, body *bytes.Buffer) {
	var resultAccount db.Account

	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	err = json.Unmarshal(data, &resultAccount)

	assert.NoError(t, err)

	assert.Equal(t, account, resultAccount)

}
