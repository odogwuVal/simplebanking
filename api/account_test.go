package api

import (
	"fmt"
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
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mockdb.NewMockStore(controller)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acct.ID)).Times(1).Return(acct, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/account/%d", acct.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// check response
	require.Equal(t, http.StatusOK, recorder.Code)

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
	}

}
