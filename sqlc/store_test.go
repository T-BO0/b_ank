package db

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(db)
	var n int64
	n = 5

	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)

	type ExecutionResult struct {
		Result *ExecuteTransferResult
		Err    error
	}
	resultsChan := make(chan ExecutionResult)
	args := ExecuteTransfereParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        decimal.New(20000000, -6),
	}

	for range n {
		go func() {
			result, err := store.ExecuteTransfer(context.Background(), args)

			resultsChan <- ExecutionResult{Result: result, Err: err}
		}()
	}

	for range n {
		exResult := <-resultsChan

		require.NoError(t, exResult.Err)
		require.NotEmpty(t, exResult.Result)

		result := exResult.Result

		require.NotEmpty(t, result.ExecutedTransfer)
		require.NotEmpty(t, result.FromAccountyEntry)
		require.NotEmpty(t, result.ToAccountEntry)
		require.NotEmpty(t, result.FromAccount)
		require.NotEmpty(t, result.ToAccount)

		var asssertionResult bool
		assertion := assert.New(t)

		asssertionResult = assertion.True(result.FromAccountyEntry.Amount.Equal(args.Amount.Neg()))
		asssertionResult = assertion.True(result.ToAccountEntry.Amount.Equal(args.Amount))
		asssertionResult = assertion.Equal(result.FromAccountyEntry.AccountID, args.FromAccountID)
		asssertionResult = assertion.Equal(result.ToAccountEntry.AccountID, args.ToAccountID)

		asssertionResult = assertion.True(result.ExecutedTransfer.Amount.Equal(args.Amount))
		asssertionResult = assertion.Equal(result.ExecutedTransfer.FromAccountID, args.FromAccountID)
		asssertionResult = assertion.Equal(result.ExecutedTransfer.ToAccountID, args.ToAccountID)

		asssertionResult = assertion.Equal(result.FromAccount.ID, args.FromAccountID)
		asssertionResult = assertion.Equal(result.ToAccount.ID, args.ToAccountID)

		diffFromAccBalance := fromAccount.Balance.Sub(result.FromAccount.Balance)
		mod := diffFromAccBalance.Mod(decimal.NewFromInt(n))
		asssertionResult = assertion.True(mod.Equal(decimal.Zero))

		diffToAccBalance := result.ToAccount.Balance.Sub(toAccount.Balance)
		mod2 := diffToAccBalance.Mod(decimal.NewFromInt(n))
		asssertionResult = assertion.True(mod2.Equal(decimal.Zero))

		require.True(t, asssertionResult)
	}
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(db)
	var n int64
	n = 10

	acc1 := createTestAccount(t)
	acc2 := createTestAccount(t)

	errorChan := make(chan error)

	for i := 0; i < int(n); i++ {
		args := ExecuteTransfereParams{
			FromAccountID: acc1.ID,
			ToAccountID:   acc2.ID,
			Amount:        decimal.New(20000000, -6),
		}

		if i%2 == 0 {
			args.FromAccountID = acc2.ID
			args.ToAccountID = acc1.ID
		}

		go func() {
			_, err := store.ExecuteTransfer(context.Background(), args)

			errorChan <- err
		}()
	}

	for range n {
		err := <-errorChan
		require.NoError(t, err)
	}
	updatedAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
