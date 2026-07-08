package repositories

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var issuer = os.Getenv("ISSUER")
var key = []byte(os.Getenv("JWT_KEY"))
var signingMethod = jwt.SigningMethodHS256

type JwtClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
type JwtManager struct {
}

func NewJwtManager() JwtManager {
	return JwtManager{}
}

func (j *JwtManager) Parse(token string) (*JwtClaims, error) {
	parser := jwt.NewParser()
	claims := &JwtClaims{}
	jwtToken, err := parser.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {

		if t.Method != signingMethod {
			return nil, fmt.Errorf("неподдерживаемый метод подписи: %v", t.Header["alg"])
		}

		return key, nil
	})
	if err != nil || !jwtToken.Valid {
		return nil, fmt.Errorf("не удалось распарсить токен: %w", err)
	}

	return claims, nil
}

func (j *JwtManager) IsAccessToken(claims *JwtClaims) bool {
	return claims.TokenType == "access"
}

func (j *JwtManager) GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	now := time.Now()

	accessToken := jwt.NewWithClaims(signingMethod,
		JwtClaims{
			TokenType: "access",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID.String(),
				Issuer:    issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			},
		})

	signedAccessToken, err := accessToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("не получилось подписать access token %w", err)
	}

	refreshToken := jwt.NewWithClaims(signingMethod,
		JwtClaims{
			TokenType: "refresh",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID.String(),
				Issuer:    issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 30)),
			},
		})

	signedRefreshToken, err := refreshToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("не получилось подписать refresh token %w", err)
	}

	return &TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (j *JwtManager) GenerateAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now()

	accessToken := jwt.NewWithClaims(signingMethod,
		JwtClaims{
			TokenType: "access",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID.String(),
				Issuer:    issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			},
		})

	signedAccessToken, err := accessToken.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("не получилось подписать access token %w", err)
	}

	return signedAccessToken, nil
}
