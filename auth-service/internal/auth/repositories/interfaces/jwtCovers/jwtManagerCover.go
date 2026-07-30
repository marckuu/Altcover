package jwtCovers

import (
	"auth-service/core/enums"
	"auth-service/internal/auth/repositories"
	"auth-service/internal/auth/repositories/interfaces"

	"github.com/google/uuid"
)

type JwtManagerCover struct {
	jwtManager repositories.JwtManager
}

func NewJwtManagerCover(jwtManager repositories.JwtManager) *JwtManagerCover {
	return &JwtManagerCover{
		jwtManager: jwtManager,
	}
}

func (j *JwtManagerCover) Parse(token string) (*interfaces.CustomClaims, error) {
	jwtClaims, err := j.jwtManager.Parse(token)
	if err != nil {
		return nil, err
	}

	registerClaims := interfaces.RegisteredClaims{
		Issuer:    jwtClaims.Issuer,
		Subject:   jwtClaims.Subject,
		ExpiresAt: jwtClaims.ExpiresAt.Time,
	}

	return &interfaces.CustomClaims{
		TokenType:        jwtClaims.TokenType,
		RegisteredClaims: registerClaims,
		Role:             jwtClaims.Role,
	}, nil
}

func (j *JwtManagerCover) GenerateTokenPair(userID uuid.UUID, userRole enums.Role) (*interfaces.TokenPair, error) {
	tokenPair, err := j.jwtManager.GenerateTokenPair(userID, userRole)
	if err != nil {
		return nil, err
	}
	return (*interfaces.TokenPair)(tokenPair), nil

}

func (j *JwtManagerCover) GenerateAccessToken(userID uuid.UUID, userRole enums.Role) (string, error) {
	return j.jwtManager.GenerateAccessToken(userID, userRole)
}

func (j *JwtManagerCover) IsAccessToken(claims *interfaces.CustomClaims) bool {
	return claims.TokenType == "access"
}
