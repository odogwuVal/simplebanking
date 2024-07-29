package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/odogwuVal/simplebanking/db/mock"
	db "github.com/odogwuVal/simplebanking/db/sqlc"
	"github.com/odogwuVal/simplebanking/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	acct := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: acct.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acct.ID)).Times(1).Return(acct, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, acct)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", acct.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotaccount db.Account
	err = json.Unmarshal(data, &gotaccount)
	require.NoError(t, err)
	require.Equal(t, account, gotaccount)
}
