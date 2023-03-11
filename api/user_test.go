package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/crackz/simple-bank/db/mock"
	db "github.com/crackz/simple-bank/db/sqlc"
	"github.com/crackz/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// func TestGetAccountAPI(t *testing.T) {
// 	inMemoryAccount := randomInMemoryAccountAccount()

// 	testCases := []struct {
// 		name          string
// 		accountID     int64
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{

// 		{
// 			name:      "OK",
// 			accountID: inMemoryAccount.ID,
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(inMemoryAccount.ID)).
// 					Times(1).
// 					Return(inMemoryAccount, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchAccount(t, recorder.Body, inMemoryAccount)
// 			},
// 		},
// 		{
// 			name:      "Invalid ID",
// 			accountID: 0,
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name:      "NOT_FOUND",
// 			accountID: inMemoryAccount.ID,
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(inMemoryAccount.ID)).
// 					Times(1).
// 					Return(db.Account{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name:      "Internal Server Error",
// 			accountID: inMemoryAccount.ID,
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					GetAccount(gomock.Any(), gomock.Eq(inMemoryAccount.ID)).
// 					Times(1).
// 					Return(db.Account{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := NewServer(store)
// 			recorder := httptest.NewRecorder()

// 			url := fmt.Sprintf("/accounts/%d", tc.accountID)
// 			req, err := http.NewRequest(http.MethodGet, url, nil)
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, req)
// 			tc.checkResponse(t, recorder)
// 		})
// 	}

// }

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}
func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v with password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	inMemoryUser, password := randomInMemoryUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

		{
			name: "Created",
			body: gin.H{
				"fullName": inMemoryUser.FullName,
				"username": inMemoryUser.Username,
				"password": password,
				"email":    inMemoryUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					FullName: inMemoryUser.FullName,
					Username: inMemoryUser.Username,
					Email:    inMemoryUser.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(inMemoryUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, inMemoryUser)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"fullName": inMemoryUser.FullName,
				"username": inMemoryUser.Username,
				"password": password,
				"email":    inMemoryUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "Duplicated Username or Email",
			body: gin.H{
				"fullName": inMemoryUser.FullName,
				"username": inMemoryUser.Username,
				"password": password,
				"email":    inMemoryUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"fullName": inMemoryUser.FullName,
				"username": inMemoryUser.Username,
				"password": password,
				"email":    "xxx",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "Too Short Password",
			body: gin.H{
				"fullName": inMemoryUser.FullName,
				"username": inMemoryUser.Username,
				"password": "123",
				"email":    inMemoryUser.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/register"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomInMemoryUser(t *testing.T) (db.User, string) {

	password := util.RandString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		Username:       util.RandOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandString(6),
		Email:          util.RandomEmail(),
	}

	return user, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
