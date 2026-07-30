package interfaces

import (
	"auth-service/core/enums"
	"time"

	"github.com/google/uuid"
)

type TokenManager interface {
	Parse(token string) (*CustomClaims, error)
	IsAccessToken(claims *CustomClaims) bool
	GenerateTokenPair(userID uuid.UUID, userRole enums.Role) (*TokenPair, error)
	GenerateAccessToken(userID uuid.UUID, userRole enums.Role) (string, error)
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type CustomClaims struct {
	TokenType        string           `json:"token_type"`
	Role             enums.Role       `json:"role"`
	RegisteredClaims RegisteredClaims `json:"registered_claims"`
}

type RegisteredClaims struct {
	Issuer    string
	Subject   string
	ExpiresAt time.Time
}
