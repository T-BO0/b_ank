package db

import (
	"context"
	"tbo-go-bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	createTestEntry(t)
}

func TestDeleteEntry(t *testing.T) {
	entry := createTestEntry(t)

	err := testQueries.DeleteEntrie(context.Background(), entry.ID)
	require.NoError(t, err)
}

func TestGetEntry(t *testing.T) {
	entry1 := createTestEntry(t)

	enrty2, err := testQueries.GetEntrie(context.Background(), entry1.ID)
	require.NoError(t, err)

	require.Equal(t, entry1.AccountID, enrty2.AccountID)
	require.Equal(t, entry1.Amount, enrty2.Amount)
	require.Equal(t, entry1.ID, enrty2.ID)
	require.Equal(t, entry1.CreatedAt, enrty2.CreatedAt)
}

func TestGetEntries(t *testing.T) {
	for range 20 {
		_ = createTestEntry(t)
	}

	args := GetEntriesParams{
		Limit:  20,
		Offset: 5,
	}
	entries, err := testQueries.GetEntries(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, v := range entries {
		require.NotEmpty(t, v)
	}

	require.LessOrEqual(t, len(entries), int(args.Limit))
}

func createTestEntry(t *testing.T) Entry {
	acc := createTestAccount(t)
	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomMonetaryAmount(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.Amount, args.Amount)
	require.Equal(t, entry.AccountID, args.AccountID)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry
}
