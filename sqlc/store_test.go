package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTX(t *testing.T) {
	store := NewStore(testDB)

	from_account := CreateRandomAccount(t)
	to_account := CreateRandomAccount(t)

	fmt.Println(">> before:", from_account.Balance, to_account.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: from_account.ID,
				ToAccountID:   to_account.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// validate transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, from_account.ID, transfer.FromAccountID)
		require.Equal(t, to_account.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// validate entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, -amount, fromEntry.Amount)
		require.Equal(t, from_account.ID, fromEntry.AccountID)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntries(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, amount, toEntry.Amount)
		require.Equal(t, to_account.ID, toEntry.AccountID)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntries(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, from_account.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, to_account.ID)

		// check account balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := from_account.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - to_account.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1*amount, 2*amount, 3*amount..., n*amount

		k := (diff1 / amount)
		require.True(t, k >= 1 && k <= int64(n))
		require.NotContains(t, existed, k)
		existed[int(k)] = true
	}
	// check the final updated account
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), from_account.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), to_account.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, from_account.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, to_account.Balance+int64(n)*amount, updatedToAccount.Balance)
}
