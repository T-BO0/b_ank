package db

import (
	"context"
	"tbo-go-bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	createTestTransfer(t)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createTestTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createTestTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
}

func TestGetTransfers(t *testing.T) {
	for range 20 {
		_ = createTestTransfer(t)
	}

	args := GetTransfersParams{
		Limit:  20,
		Offset: 5,
	}

	transfers, err := testQueries.GetTransfers(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	require.LessOrEqual(t, len(transfers), int(args.Limit))
	for _, v := range transfers {
		require.NotEmpty(t, v)
	}
}

func createTestTransfer(t *testing.T) Transfer {
	accFrom := createTestAccount(t)
	accTo := createTestAccount(t)

	args := CreateTransferParams{
		FromAccountID: accFrom.ID,
		ToAccountID:   accTo.ID,
		Amount:        util.RandomMonetaryAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, args.Amount, transfer.Amount)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}
