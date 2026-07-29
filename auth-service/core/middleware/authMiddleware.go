package middleware

import (
	"auth-service/core/enums"
	errors2 "auth-service/core/errors"
	"auth-service/internal/auth/repositories/interfaces"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const ClaimsKey = "claims"

type AuthMiddleware struct {
	tokenManager interfaces.TokenManager
}

func NewAuthMiddleware(tokenManager interfaces.TokenManager) AuthMiddleware {
	return AuthMiddleware{
		tokenManager: tokenManager,
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	claims := ctx.Value(ClaimsKey).(*interfaces.CustomClaims)
	userID, err := uuid.Parse(claims.RegisteredClaims.Subject)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

func (a *AuthMiddleware) ProcessToken(header string) (*interfaces.CustomClaims, error) {
	// Получить access токен из запроса
	headerParts := strings.Split(header, " ")
	if headerParts[0] != "Bearer" {
		fmt.Println("Не найден Bearer в заголовке")
		return nil, errors.New("неподдерживаемый вид авторизации")
	}

	// Проверка что токен был подписан этой системой и никто его не поменял
	claims, err := a.tokenManager.Parse(headerParts[1])
	if err != nil {
		fmt.Println("ошибка парсинга access токена")
		return nil, err
	}

	// Проверка что это именно access токен
	if isAccess := a.tokenManager.IsAccessToken(claims); !isAccess {
		fmt.Println("переданный токен не является access токеном")
		return nil, errors.New("токен не является access токеном")
	}

	return claims, nil
}

func (a *AuthMiddleware) RequireRole(role enums.Role, handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		claims, err := a.ProcessToken(authHeader)
		if err != nil {
			errors2.SendErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		if claims.Role < role {
			errors2.SendErrorResponse(w, errors.New("недостаточно прав для выполнения запроса"), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}
