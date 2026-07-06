package interfaces

import (
	"time"

	"github.com/google/uuid"
)

type TokenManager interface {
	Parse(token string) (*CustomClaims, error)
	IsAccessToken(claims *CustomClaims) bool
	GenerateTokenPair(userID uuid.UUID) (*TokenPair, error)
	GenerateAccessToken(userID uuid.UUID) (string, error)
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type CustomClaims struct {
	TokenType        string `json:"token_type"`
	RegisteredClaims RegisteredClaims
}

type RegisteredClaims struct {
	Issuer    string
	Subject   string
	ExpiresAt time.Time
}
