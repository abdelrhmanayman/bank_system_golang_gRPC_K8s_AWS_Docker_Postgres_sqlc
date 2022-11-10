package api

import (
	mockdb "banksystem/db/mock"
	db "banksystem/db/sqlc"
	"banksystem/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type loginTestCase struct {
	name              string
	payload           LoginRequest
	testStubs         func(store *mockdb.MockStore)
	checkTestResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := util.GenerateRandomString(8)
	hashedPassword, err := util.HashPassword(password)

	assert.NoError(t, err)

	return db.User{
		Username:  util.CreateRandomOwner(),
		HashedPwd: hashedPassword,
		Email:     util.CreateRandomEmail(),
	}, password
}

func TestLoginController(t *testing.T) {
	assert := assert.New(t)
	dummyUser, password := createRandomUser(t)
	nonExistingUser := util.CreateRandomOwner()

	testCases := []loginTestCase{
		{
			name: "Ok",
			payload: LoginRequest{
				Username: dummyUser.Username,
				Password: password,
			},
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(dummyUser.Username)).Times(1).Return(dummyUser, nil)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Non Existing User",
			payload: LoginRequest{
				Username: nonExistingUser,
				Password: password,
			},
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(nonExistingUser)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Wrong Password",
			payload: LoginRequest{
				Username: dummyUser.Username,
				Password: util.GenerateRandomString(8),
			},
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(dummyUser.Username)).Times(1).Return(dummyUser, nil)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			payload: LoginRequest{
				Password: util.GenerateRandomString(8),
			},
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(dummyUser.Username)).Times(0)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)
		store := mockdb.NewMockStore(ctrl)

		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		testCase.testStubs(store)

		var bodyBuffer bytes.Buffer
		json.NewEncoder(&bodyBuffer).Encode(testCase.payload)

		request, err := http.NewRequest(http.MethodPost, "/auth/login", &bodyBuffer)
		assert.NoError(err)

		server.router.ServeHTTP(recorder, request)
		testCase.checkTestResponse(t, recorder)
	}

}
