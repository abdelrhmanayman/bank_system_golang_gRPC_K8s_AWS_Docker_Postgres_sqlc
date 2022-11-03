package api

import (
	mockdb "banksystem/db/mock"
	db "banksystem/db/sqlc"
	"banksystem/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAccountController(t *testing.T) {
	assert := assert.New(t)
	account := createDummyAccount()

	testCases := []struct {
		name              string
		accountID         int64
		testStubs         func(store *mockdb.MockStore)
		checkTestResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusOK, recorder.Code)
				matchAccountBodyCheck(t, account, recorder.Body)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), 0).Times(0)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)

			server := SetupRoutes(store)
			recorder := httptest.NewRecorder()
			testCase.testStubs(store)

			url := fmt.Sprintf("/accounts/%d", testCase.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkTestResponse(t, recorder)

		})
	}

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

	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	err = json.Unmarshal(data, &resultAccount)

	assert.NoError(t, err)

	assert.Equal(t, account, resultAccount)

}
