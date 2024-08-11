package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/thaian1234/simplebank/utils"
)

func TestNewJWTMakerV5(t *testing.T) {
	validKey := utils.RandomString(32)
	invalidKey := utils.RandomString(31)

	maker, err := NewJWTMakerV5(validKey)
	require.NoError(t, err)
	require.NotNil(t, maker)

	maker, err = NewJWTMakerV5(invalidKey)
	require.Error(t, err)
	require.Nil(t, maker)
}

func TestJWTMakerV5_GenerateToken(t *testing.T) {
	maker, err := NewJWTMakerV5(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	token, err := maker.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestJWTMakerV5_VerifyToken(t *testing.T) {
	maker, err := NewJWTMakerV5(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	token, err := maker.GenerateToken(username, duration)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.Equal(t, username, payload.Username)
}

func TestJWTMakerV5_ExpiredToken(t *testing.T) {
	maker, err := NewJWTMakerV5(utils.RandomString(32))
	require.NoError(t, err)

	token, err := maker.GenerateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, jwt.ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestJWTMakerV5_InvalidToken(t *testing.T) {
	maker, err := NewJWTMakerV5(utils.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("invalid.token.here")
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewUserClaims(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMakerV5(utils.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, jwt.ErrSignatureInvalid.Error())
	require.Nil(t, payload)
}
