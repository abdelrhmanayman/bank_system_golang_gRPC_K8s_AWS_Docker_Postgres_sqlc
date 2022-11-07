package api

import (
	mockdb "banksystem/db/mock"
	db "banksystem/db/sqlc"
	"banksystem/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type userTestCase struct {
	name              string
	payload           CreateUserRequest
	testStubs         func(store *mockdb.MockStore)
	checkTestResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

func createDummyUser() CreateUserRequest {

	return CreateUserRequest{
		Username: util.CreateRandomOwner(),
		Password: util.GenerateRandomString(8),
		Email:    util.CreateRandomEmail(),
	}
}

func TestCreateUserControllerTest(t *testing.T) {
	assert := assert.New(t)
	dummyUser := createDummyUser()

	testCases := []userTestCase{
		{
			name:    "OK",
			payload: dummyUser,
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{
					Username: dummyUser.Username,
					Email:    dummyUser.Email,
				}, nil)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusOK, recorder.Code)
				matchUserBody(t, dummyUser, recorder.Body)
			},
		},
		{
			name:    "Forbidden",
			payload: dummyUser,
			testStubs: func(store *mockdb.MockStore) {
				var vioError pq.ErrorCode = "23505"
				mockedConstraintViolation := pq.Error{
					Code: vioError,
				}
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, &mockedConstraintViolation)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:    "Internal Server Error",
			payload: dummyUser,
			testStubs: func(store *mockdb.MockStore) {

				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			payload: CreateUserRequest{
				Username: "kasdf",
				Password: "sdfsdf",
				Email:    "sdfsadf",
			},
			testStubs: func(store *mockdb.MockStore) {

				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkTestResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)
		server := SetupRoutes(store)
		recorder := httptest.NewRecorder()

		testCase.testStubs(store)
		var bodyBuffer bytes.Buffer
		json.NewEncoder(&bodyBuffer).Encode(testCase.payload)

		request, err := http.NewRequest(http.MethodPost, "/users", &bodyBuffer)
		assert.NoError(err)

		server.router.ServeHTTP(recorder, request)
		testCase.checkTestResponse(t, recorder)
	}

}

func matchUserBody(t *testing.T, user CreateUserRequest, body *bytes.Buffer) {
	var responseBody db.User

	bodyBytes, err := io.ReadAll(body)

	assert.NoError(t, err)

	err = json.Unmarshal(bodyBytes, &responseBody)

	assert.NoError(t, err)

	assert.Equal(t, user.Username, responseBody.Username)
	assert.Equal(t, user.Email, responseBody.Email)

}
