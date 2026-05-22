package handlers

import (
	"auth-service/domains"
	"auth-service/dto"
	"auth-service/enums"
	"auth-service/handlers/tools"
	"auth-service/repositories"
	"auth-service/services"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type HTTPAuthHandlers struct {
	tokenService services.TokenService
	jwtManager   repositories.JWTManager
	userService  services.UserService
	ctx          context.Context
}

func NewHTTPAuthHandler() HTTPAuthHandlers {
	return HTTPAuthHandlers{
		tokenService: services.TokenService{},
		jwtManager:   repositories.JWTManager{},
		userService:  services.UserService{},
		ctx:          nil,
	}
}

func (a *HTTPAuthHandlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Получить ник и пароль из запроса
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		fmt.Printf("не удалось получить ник и пароль из запроса, %v", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Попробовать найти пользователя в системе по нику
	_, err := a.userService.GetUserByNickname(a.ctx, loginRequest.Nickname)
	if err == nil {
		fmt.Printf("пользователь уже существует, %v", err)
		tools.SendErrorResponse(w, errors.New("пользователь уже существует"), http.StatusBadRequest)
		return
	}

	// Получить хэш пароля
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), 10)
	if err != nil {
		fmt.Printf("ошибка получения хэша пароля, %v", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Создать нового пользователя
	var user domains.User

	user.Role = enums.User
	user.Nickname = loginRequest.Nickname
	user.PasswordHash = PasswordHash

	// Добавить пользователя
	if err = a.userService.AddUser(a.ctx, user); err != nil {
		fmt.Printf("ошибка добавления нового пользователя, %v", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Отправить ответ об успешной регистраци
}

func (a *HTTPAuthHandlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Получить ник и пароль из запроса
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		fmt.Printf("не удалось получить ник и пароль из запроса, %v", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Попробовать найти пользователя в системе по нику
	user, err := a.userService.GetUserByNickname(a.ctx, loginRequest.Nickname)
	if err != nil {
		fmt.Printf("пользователь не найден, %v", err)
		tools.SendErrorResponse(w, errors.New("пользователь не найден"), http.StatusBadRequest)
		return
	}

	// Сравнить хэши паролей
	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(loginRequest.Password)); err != nil {
		fmt.Printf("неверный пароль, %v", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Сгенерировать пару токенов
	tokenPair, err := a.jwtManager.GenerateTokenPair(user.ID)
	if err != nil {
		fmt.Printf("не удалось сгенерировать jwt токены, %v", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Записать refresh токен в бд
	if err = a.tokenService.AddRefreshToken(a.ctx, tokenPair.RefreshToken); err != nil {
		fmt.Printf("не удалось сохранить refresh токен, %v", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Записать refresh в cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "auth/refresh",
		MaxAge:   60 * 60 * 24 * 30,
	})

	// access отправить в ответе
	err = json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenPair.AccessToken,
	})
	if err != nil {
		fmt.Printf("не удалось отправить access токен в ответе, %v", err)
	}

}

func (a *HTTPAuthHandlers) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// Получить refresh из cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Printf("отсутствует refresh токен, %v", err)
		tools.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	// Проверить валидность refresh токена
	claims, err := a.jwtManager.Parse(cookie.Value)
	if err != nil {
		fmt.Printf("неккоректный refresh токен, %v", err)
		tools.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	if err = a.jwtManager.Validate(claims); err != nil {
		fmt.Printf("невалидный refresh токен, %v", err)
		tools.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	// Проверить в бд, что токен не отозван
	isRevoked, err := a.tokenService.IsTokenRevoked(a.ctx, cookie.Value)
	if err != nil {
		fmt.Printf("не удалось проверить токен, %v", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if isRevoked {
		fmt.Printf("переданный refresh токен отозван, %v", err)
		tools.SendErrorResponse(w, errors.New("переданный refresh токен отозван"), http.StatusInternalServerError)
		return
	}

	// Сгенерировать новый access токен и вернуть его в ответе
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		fmt.Printf("ошибка получения id пользователя из refresh токена, %v", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	accessToken, err := a.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		fmt.Printf("не удалось сгенерировать access токен, %v", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
	if err != nil {
		fmt.Printf("не удалось отправить access токен, %v", err)
		// Логировать ошибку
	}
}

func (a *HTTPAuthHandlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Получить refresh из cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Printf("отсутствует refresh токен, %v", err)
		tools.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	// Проверить валидность refresh токена
	_, err = a.jwtManager.Parse(cookie.Value)
	if err != nil {
		fmt.Printf("неккоректный refresh токен, %v", err)
		tools.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	// Отозвать refresh токен из бд
	if err = a.tokenService.DeleteRefreshToken(a.ctx, cookie.Value); err != nil {
		fmt.Printf("ошибка удаления refresh токена из базы, %v", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Удалить refresh из cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "auth/refresh",
		MaxAge:   -1,
	})
}
