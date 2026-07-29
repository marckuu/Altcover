package repositories

import (
	enums2 "book-cover-service/core/enums"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var issuer = os.Getenv("ISSUER")
var key = []byte(os.Getenv("JWT_KEY"))
var signingMethod = jwt.SigningMethodHS256

type JWTManager struct {
}

func NewJWTManager() JWTManager {
	return JWTManager{}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
	Role enums2.Role
}

// Метод ParseWithClaims получает возвращаемый из аноним функции ключ, высчитывает подпись и проверяет
// валидна ли она
func (j *JWTManager) Parse(token string) (*CustomClaims, error) {
	parser := jwt.NewParser()
	claims := &CustomClaims{}
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

func (j *JWTManager) Validate(claims *CustomClaims) error {
	if claims.ExpiresAt.Compare(time.Now()) == 1 {
		return errors.New("истёк срок действия токена")
	}
	return nil
}

func (j *JWTManager) IsAccessToken(claims *CustomClaims) bool {
	return claims.TokenType == "access"
}

func (j *JWTManager) GenerateTokenPair(userID uuid.UUID, userRole enums2.Role) (*TokenPair, error) {
	now := time.Now()

	accessToken := jwt.NewWithClaims(signingMethod,
		CustomClaims{
			TokenType: "access",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID.String(),
				Issuer:    issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			},
			Role: userRole,
		})

	signedAccessToken, err := accessToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("не получилось подписать access token %w", err)
	}

	refreshToken := jwt.NewWithClaims(signingMethod,
		CustomClaims{
			TokenType: "refresh",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID.String(),
				Issuer:    issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 30)),
			},
			Role: userRole,
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

func (j *JWTManager) GenerateAccessToken(userID uuid.UUID, userRole enums2.Role) (string, error) {
	now := time.Now()

	accessToken := jwt.NewWithClaims(signingMethod,
		CustomClaims{
			TokenType: "access",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   userID.String(),
				Issuer:    issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			},
			Role: userRole,
		})

	signedAccessToken, err := accessToken.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("не получилось подписать access token %w", err)
	}

	return signedAccessToken, nil
}
