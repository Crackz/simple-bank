package token

import (
	"testing"
	"time"

	"github.com/crackz/simple-bank/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandString(32))
	require.NoError(t, err)

	username := util.RandOwner()
	duration := time.Hour * 24
	issueAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt.Time, issueAt, time.Second)
	require.WithinDuration(t, payload.ExpiresAt.Time, expiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandString(32))
	require.NoError(t, err)

	username := util.RandOwner()

	token, err := maker.CreateToken(username, -time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)
	require.Nil(t, payload)
}

func TestForbidNoneSignatureTypeJWTToken(t *testing.T) {
	payload, err := NewPayload(util.RandOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenMalformed)
	require.Nil(t, payload)
}
