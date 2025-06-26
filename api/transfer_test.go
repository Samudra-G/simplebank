package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Samudra-G/simplebank/db/mock"
	db "github.com/Samudra-G/simplebank/db/sqlc"
	"github.com/Samudra-G/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T){
	account1 := randomAccount()
	account2 := randomAccount()
	account1.Currency = "USD"
	account2.Currency = "USD"

	amount := int64(10)

	transferResult := db.TransferTxResult{
		Transfer: db.Transfer{
			ID: 		   util.RandomInt(1, 1000),
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount: 	   amount,	
		},
		FromAccount: account1,
		ToAccount:   account2,
		FromEntry: db.Entry{
			ID:        util.RandomInt(1, 1000),
			AccountID: account1.ID,
			Amount:    -amount,
		},
		ToEntry: db.Entry{
			ID:        util.RandomInt(1, 1000),
			AccountID: account2.ID,
			Amount:    amount,
		},
	}

	testCases := []struct {
		name 		  string
		body 		  gin.H 
		buildStubs	  func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(transferResult, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchTransferTxResult(t, rec.Body, transferResult)
			},
		},
		{
			name: "CurrencyMismatch",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "EUR",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account_id": 99999,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(int64(99999))).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			name: "TransferTxFails",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
		{
			name: "InvalidInput",
			body: gin.H{
				"from_account_id": 0,
				"to_account_id":   account2.ID,
				"amount":          -10,
				"currency":        "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}
	// Run the tests
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStubs(mockStore)

			server := NewServer(mockStore)
			rec := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func requireBodyMatchTransferTxResult(t *testing.T, body *bytes.Buffer, expected db.TransferTxResult) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got db.TransferTxResult
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	require.Equal(t, expected.Transfer.FromAccountID, got.Transfer.FromAccountID)
	require.Equal(t, expected.Transfer.ToAccountID, got.Transfer.ToAccountID)
	require.Equal(t, expected.Transfer.Amount, got.Transfer.Amount)

	require.Equal(t, expected.FromAccount.ID, got.FromAccount.ID)
	require.Equal(t, expected.ToAccount.ID, got.ToAccount.ID)

	require.Equal(t, expected.FromEntry.Amount, got.FromEntry.Amount)
	require.Equal(t, expected.ToEntry.Amount, got.ToEntry.Amount)
}