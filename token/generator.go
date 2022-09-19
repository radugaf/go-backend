package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/golang-jwt/jwt/v4"
	"github.com/o1egl/paseto"
)

const minSecretKeySize = 32

// JWTGenerator is a JSON Web Token generator
type JWTGenerator struct {
	secretKey string
}

// NewJWTGenerator creates a new JWTGenerator
func NewJWTGenerator(secretKey string) (Token, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTGenerator{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (generator *JWTGenerator) GenerateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(generator.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (generator *JWTGenerator) ValidateToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(generator.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

type PasetoGenerator struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoGenerator creates a new PasetoGenerator
func NewPasetoGenerator(symmetricKey string) (Token, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	generator := &PasetoGenerator{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return generator, nil
}

// CreateTokegeneratorates a new token for a specific username and duration
func (generator *PasetoGenerator) GenerateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := generator.paseto.Encrypt(generator.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (generator *PasetoGenerator) ValidateToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := generator.paseto.Decrypt(token, generator.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
