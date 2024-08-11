package token

import (
	"github.com/stretchr/testify/require"
	"github.com/thaian1234/simplebank/utils"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	tokenString := utils.RandomString(32)
	maker, err := NewPasetoMaker([]byte(tokenString))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Hour

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}
