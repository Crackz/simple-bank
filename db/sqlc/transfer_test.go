package db

import (
	"context"
	"testing"

	"github.com/crackz/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, from_account *Account, to_account *Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Amount:        util.RandomInt(1, 10000),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)

	createRandomTransfer(t, &from_account, &to_account)
}

func TestGetTransfer(t *testing.T) {
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)

	transfer := createRandomTransfer(t, &from_account, &to_account)

	foundTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.Equal(t, transfer.ID, foundTransfer.ID)
	require.Equal(t, transfer.FromAccountID, foundTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, foundTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, foundTransfer.Amount)
	require.Equal(t, transfer.CreatedAt, foundTransfer.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	from_account := createRandomAccount(t)
	to_account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, &from_account, &to_account)
	}

	arg := ListTransfersParams{
		FromAccountID: from_account.ID,
		ToAccountID:   to_account.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}
