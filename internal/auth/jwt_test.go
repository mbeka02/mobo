package auth

import (
	"testing"
	"time"

	"github.com/mbeka02/ticketing-service/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandString(32))
	require.NoError(t, err)
	email := utils.RandEmail()
	userId := utils.RandUUID()
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := time.Now().Add(duration)

	token, err := maker.Create(userId, email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := maker.Verify(token)
	require.NoError(t, err)
	require.NotEmpty(t, claims)
	require.Equal(t, email, claims.Email)
	require.Equal(t, userId, claims.UserID)
	require.WithinDuration(t, issuedAt, claims.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, claims.ExpiresAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandString(32))
	require.NoError(t, err)
	email := utils.RandEmail()
	userId := utils.RandUUID()
	duration := -time.Minute
	token, err := maker.Create(userId, email, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := maker.Verify(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, claims)
}
