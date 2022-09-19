package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/radugaf/simplebank/tools"
	"github.com/stretchr/testify/require"
)

func TestJWTGenerator(t *testing.T) {
	generator, err := NewJWTGenerator(tools.RandomString(32))
	require.NoError(t, err)

	username := tools.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := generator.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = generator.ValidateToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	generator, err := NewJWTGenerator(tools.RandomString(32))
	require.NoError(t, err)

	token, payload, err := generator.GenerateToken(tools.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = generator.ValidateToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgIsNone(t *testing.T) {
	payload, err := NewPayload(tools.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	generator, err := NewJWTGenerator(tools.RandomString(32))
	require.NoError(t, err)

	payload, err = generator.ValidateToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestPasetoGenerator(t *testing.T) {
	generator, err := NewPasetoGenerator(tools.RandomString(32))
	require.NoError(t, err)

	username := tools.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := generator.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = generator.ValidateToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	generator, err := NewPasetoGenerator(tools.RandomString(32))
	require.NoError(t, err)

	token, payload, err := generator.GenerateToken(tools.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = generator.ValidateToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
