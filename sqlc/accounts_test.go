package db

import (
	"context"
	"tbo-go-bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestDeleteAccount(t *testing.T) {
	acc := createTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc.ID)
	require.NoError(t, err)
}

func TestGetAccount(t *testing.T) {
	acc := createTestAccount(t)

	getAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAcc)

	require.Equal(t, acc.ID, getAcc.ID)
	require.Equal(t, acc.Currency, getAcc.Currency)
	require.Equal(t, acc.OwnerName, getAcc.OwnerName)
	require.Equal(t, acc.Balance, getAcc.Balance)
	require.Equal(t, acc.CreatedAt, getAcc.CreatedAt)
}

func TestGetAccounts(t *testing.T) {
	for range 20 {
		_ = createTestAccount(t)
	}

	args := GetAccountsParams{
		Limit:  20,
		Offset: 5,
	}
	accs, err := testQueries.GetAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accs)

	require.LessOrEqual(t, len(accs), int(args.Limit))

	for _, v := range accs {
		require.NotEmpty(t, v)
	}
}

func TestUpdateAccount(t *testing.T) {
	acc := createTestAccount(t)

	args := UpdateAccountBalanceParams{
		ID:     acc.ID,
		Amount: util.RandomMonetaryAmount(),
	}

	updatedAcc, err := testQueries.UpdateAccountBalance(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)
	require.NotEmpty(t, updatedAcc.Balance)
	require.NotEmpty(t, updatedAcc.OwnerName)
	require.NotEmpty(t, updatedAcc.Currency)
	require.NotEmpty(t, updatedAcc.ID)
	require.NotEmpty(t, updatedAcc.CreatedAt)

	amm := (updatedAcc.Balance.Sub(acc.Balance))

	require.Equal(t, amm, args.Amount)
}
