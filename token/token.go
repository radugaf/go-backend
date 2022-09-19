package token

import "time"

// Token is an interface for managing authentication tokens
type Token interface {
	// GenerateToken generates a new token for a specific user
	GenerateToken(username string, duration time.Duration) (string, *Payload, error)
	// ValidateToken checks if the token is valid or not
	ValidateToken(token string) (*Payload, error)
}
