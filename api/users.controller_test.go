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
	"reflect"
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

type createUserMatcher struct {
	password string
	userArgs db.CreateUserParams
}

func (userMatcher createUserMatcher) Matches(x interface{}) bool {
	userArgs, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}

	isPasswordsMatches := util.CompareHashedPasswords(userMatcher.password, userArgs.HashedPwd)

	if !isPasswordsMatches {
		return false
	}

	userMatcher.userArgs.HashedPwd = userArgs.HashedPwd

	return reflect.DeepEqual(userMatcher.userArgs, userArgs)
}

func (userMatcher createUserMatcher) String() string {
	return fmt.Sprintf("user with username: %s is matched successfully", userMatcher.userArgs.Username)
}

func matchCreateUserArgs(args db.CreateUserParams, password string) gomock.Matcher {
	return createUserMatcher{
		password: password,
		userArgs: args,
	}
}

func TestCreateUserControllerTest(t *testing.T) {
	assert := assert.New(t)
	dummyUser := createDummyUser()
	createUserArgs := db.CreateUserParams{
		Username: dummyUser.Username,
		Email:    dummyUser.Email,
	}
	testCases := []userTestCase{
		{
			name:    "OK",
			payload: dummyUser,
			testStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), matchCreateUserArgs(createUserArgs, dummyUser.Password)).Times(1).Return(db.User{
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
		server := newTestServer(t, store)
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
