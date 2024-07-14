package db

import (
	"context"
	"testing"
	"time"

	"github.com/odogwuVal/simplebanking/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntries(t *testing.T, account Account) Entry {
	args := CreateEntriesParams{
		Amount:    util.RandomAmount(),
		AccountID: account.ID,
	}
	entry, err := testQueries.CreateEntries(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	require.NotZero(t, entry.AccountID)
	require.NotZero(t, entry.Amount)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	CreateRandomEntries(t, account)
}

func TestGetEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	entry := CreateRandomEntries(t, account)
	ent, err := testQueries.GetEntries(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, ent)

	require.Equal(t, entry.AccountID, ent.AccountID)
	require.Equal(t, entry.ID, ent.ID)
	require.Equal(t, entry.Amount, ent.Amount)

	require.WithinDuration(t, entry.CreatedAt, ent.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomEntries(t, account)
	}

	args := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entry, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entry, 5)

	for _, ent := range entry {
		require.NotEmpty(t, ent)
	}
}
