package db

import (
	"context"
	"testing"
	"time"

	"github.com/odogwuVal/simplebanking/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, from_account, to_account Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Amount:        util.RandomAmount(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	from_account := CreateRandomAccount(t)
	to_account := CreateRandomAccount(t)
	CreateRandomTransfer(t, from_account, to_account)
}

func TestGetTransfer(t *testing.T) {
	from_account := CreateRandomAccount(t)
	to_account := CreateRandomAccount(t)
	transfer := CreateRandomTransfer(t, from_account, to_account)

	trans, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trans)

	require.Equal(t, transfer.ID, trans.ID)
	require.Equal(t, transfer.FromAccountID, trans.FromAccountID)
	require.Equal(t, transfer.ToAccountID, trans.ToAccountID)
	require.Equal(t, transfer.Amount, trans.Amount)
	require.WithinDuration(t, transfer.CreatedAt, trans.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	from_account := CreateRandomAccount(t)
	to_account := CreateRandomAccount(t)

	args := ListTransfersParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Limit:         5,
		Offset:        0,
	}

	for i := 0; i < 5; i++ {
		CreateRandomTransfer(t, from_account, to_account)
	}

	transfer, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	for _, trans := range transfer {
		require.NotEmpty(t, trans)
		require.True(t, trans.FromAccountID == from_account.ID || trans.ToAccountID == to_account.ID)
	}
}
