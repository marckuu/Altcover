package middleware

import (
	"auth-service/db/repositories"
	"auth-service/handlers/tools"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const ClaimsKey = "claims"

type AuthMiddleware struct {
	jwtManager repositories.JWTManager
}

func NewAuthMiddleware(JWTManager repositories.JWTManager) AuthMiddleware {
	return AuthMiddleware{
		jwtManager: JWTManager,
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	claims := ctx.Value("claims").(repositories.CustomClaims)
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

func (a *AuthMiddleware) ProcessToken(header string) (*repositories.CustomClaims, error) {
	// Получить access токен из запроса
	headerParts := strings.Split(header, " ")
	if headerParts[0] != "Bearer" {
		fmt.Println("Не найден Bearer в заголовке")
		return nil, errors.New("неподдерживаемый вид авторизации")
	}

	// Проверка что токен был подписан этой системой и никто его не поменял
	claims, err := a.jwtManager.Parse(headerParts[1])
	if err != nil {
		fmt.Println("ошибка парсинга access токена")
		return nil, err
	}

	// Валидация токена
	if err = a.jwtManager.Validate(claims); err != nil {
		fmt.Println("у переданного токена истёк срок жизни")
		return nil, errors.New("передан токен с истёкшим сроком действия")
	}

	// Проверка что это именно access токен
	if isAccess := a.jwtManager.IsAccessToken(claims); !isAccess {
		fmt.Println("переданный токен не является access токеном")
		return nil, errors.New("токен не является access токеном")
	}

	return claims, nil
}

func (a *AuthMiddleware) Auth(handler http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		claims, err := a.ProcessToken(authHeader)
		if err != nil {
			tools.SendErrorResponse(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	}
}
