package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	db "github.com/radugaf/simplebank/db/sqlc"
	"github.com/radugaf/simplebank/token"
	"github.com/radugaf/simplebank/tools"
	"github.com/stretchr/testify/require"

	mockdb "github.com/radugaf/simplebank/db/mock"
)

func TestGetBankAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	bankAccount := randomAccount(user.Username)

	testCases := []struct {
		setupAuth     func(t *testing.T, request *http.Request, tokenGenerator token.Token)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		buildStubs    func(store *mockdb.MockStore)
		name          string
		bankAccountID int64
	}{
		{
			name:          "OK",
			bankAccountID: bankAccount.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetBankAccount(gomock.Any(), gomock.Eq(bankAccount.ID)).Times(1).Return(bankAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBankAccount(t, recorder.Body, bankAccount)
			},
		},
		{
			name:          "NotFound",
			bankAccountID: bankAccount.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetBankAccount(gomock.Any(), gomock.Eq(bankAccount.ID)).Times(1).Return(db.BankAccount{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:          "InternalError",
			bankAccountID: bankAccount.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetBankAccount(gomock.Any(), gomock.Eq(bankAccount.ID)).Times(1).Return(db.BankAccount{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:          "InvalidID",
			bankAccountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/bank_accounts/%d", tc.bankAccountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateBankAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		setupAuth     func(t *testing.T, request *http.Request, tokenGenerator token.Token)
		checkResponse func(recoder *httptest.ResponseRecorder)
		buildStubs    func(store *mockdb.MockStore)
		body          gin.H
		name          string
	}{
		{
			name: "OK",
			body: gin.H{"currency": account.Currency},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateBankAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}

				store.EXPECT().
					CreateBankAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBankAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBankAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBankAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.BankAccount{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"currency": "invalid",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBankAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/bank_accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListBankAccountsAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	accounts := make([]db.BankAccount, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		setupAuth     func(t *testing.T, request *http.Request, tokenGenerator token.Token)
		checkResponse func(recoder *httptest.ResponseRecorder)
		buildStubs    func(store *mockdb.MockStore)
		name          string
		query         Query
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListBankAccountsParams{
					Owner:  user.Username,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListBankAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBankAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBankAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBankAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.BankAccount{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBankAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenGenerator token.Token) {
				addAuthorization(t, request, tokenGenerator, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBankAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/bank_accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenGenerator)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomAccount(owner string) db.BankAccount {
	return db.BankAccount{
		ID:       tools.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  tools.RandomMoney(),
		Currency: tools.RandomCurrency(),
	}
}

func requireBodyMatchBankAccount(t *testing.T, body *bytes.Buffer, bankAccount db.BankAccount) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.BankAccount

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, bankAccount, gotAccount)
}

func requireBodyMatchBankAccounts(t *testing.T, body *bytes.Buffer, accounts []db.BankAccount) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.BankAccount
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
