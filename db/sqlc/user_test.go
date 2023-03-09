package db

import (
	"context"
	"testing"

	"github.com/crackz/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandString(8),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)

}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	foundUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, foundUser.Username)
	require.Equal(t, user.FullName, foundUser.FullName)
	require.Equal(t, user.HashedPassword, foundUser.HashedPassword)
	require.Equal(t, user.Email, foundUser.Email)
	require.Equal(t, user.CreatedAt, foundUser.CreatedAt)
	require.Equal(t, user.PasswordChangedAt, foundUser.PasswordChangedAt)

}
