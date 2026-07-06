package transport

import (
	"auth-service/core/domains"
	"auth-service/core/enums"
	coreErrors "auth-service/core/errors"
	logs "auth-service/core/logger"
	globInterfaces "auth-service/internal/auth/repositories/interfaces"
	servicesInterfaces "auth-service/internal/auth/services/interfaces"
	"auth-service/internal/auth/transport/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type HTTPAuthHandlers struct {
	tokenService servicesInterfaces.TokenService
	tokenManager globInterfaces.TokenManager
	userService  servicesInterfaces.UserService
	ctx          context.Context
	logger       logs.Logger
}

func NewHTTPAuthHandler(tokenService servicesInterfaces.TokenService,
	tokenManager globInterfaces.TokenManager,
	userService servicesInterfaces.UserService,
	ctx context.Context,
	logger logs.Logger) HTTPAuthHandlers {
	return HTTPAuthHandlers{
		tokenService: tokenService,
		tokenManager: tokenManager,
		userService:  userService,
		ctx:          ctx,
		logger:       logger,
	}
}

func (a *HTTPAuthHandlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		a.logger.Error(fmt.Errorf("не удалось получить ник и пароль из запроса: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	_, err := a.userService.GetUserByNickname(a.ctx, loginRequest.Nickname)
	if err == nil {
		a.logger.Error(fmt.Errorf("пользователь уже существует: %w", err).Error())
		coreErrors.SendErrorResponse(w, errors.New("пользователь уже существует"), http.StatusBadRequest)
		return
	}

	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), 10)
	if err != nil {
		a.logger.Error(fmt.Errorf("ошибка получения хэша пароля: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var user domains.User

	user.Role = enums.User
	user.Nickname = loginRequest.Nickname
	user.PasswordHash = PasswordHash

	if err = a.userService.AddUser(a.ctx, user); err != nil {
		a.logger.Error(fmt.Errorf("ошибка добавления нового пользователя: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *HTTPAuthHandlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		a.logger.Error(fmt.Errorf("не удалось получить ник и пароль из запроса: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, err := a.userService.GetUserByNickname(a.ctx, loginRequest.Nickname)
	if err != nil {
		a.logger.Error(fmt.Errorf("пользователь не найден: %w", err).Error())
		coreErrors.SendErrorResponse(w, errors.New("пользователь не найден"), http.StatusBadRequest)
		return
	}

	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(loginRequest.Password)); err != nil {
		a.logger.Error(fmt.Errorf("неверный пароль: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	tokenPair, err := a.tokenManager.GenerateTokenPair(user.ID)
	if err != nil {
		a.logger.Error(fmt.Errorf("не удалось сгенерировать jwtCovers токены: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = a.tokenService.AddRefreshToken(a.ctx, tokenPair.RefreshToken); err != nil {
		a.logger.Error(fmt.Errorf("не удалось сохранить refresh токен: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "auth/refresh",
		MaxAge:   60 * 60 * 24 * 30,
	})

	err = json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenPair.AccessToken,
	})
	if err != nil {
		a.logger.Error(fmt.Errorf("не удалось отправить access токен в ответе: %w", err).Error())
	}

}

func (a *HTTPAuthHandlers) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		a.logger.Error(fmt.Errorf("отсутствует refresh токен: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	claims, err := a.tokenManager.Parse(cookie.Value)
	if err != nil {
		a.logger.Error(fmt.Errorf("неккоректный refresh токен: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	isRevoked, err := a.tokenService.IsTokenRevoked(a.ctx, cookie.Value)
	if err != nil {
		a.logger.Error(fmt.Errorf("не удалось проверить токен: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if isRevoked {
		a.logger.Error(fmt.Errorf("переданный refresh токен отозван: %w", err).Error())
		coreErrors.SendErrorResponse(w, errors.New("переданный refresh токен отозван"), http.StatusInternalServerError)
		return
	}

	userID, err := uuid.Parse(claims.RegisteredClaims.Subject)
	if err != nil {
		a.logger.Error(fmt.Errorf("ошибка получения id пользователя из refresh токена: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	accessToken, err := a.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		a.logger.Error(fmt.Errorf("не удалось сгенерировать access токен: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
	if err != nil {
		a.logger.Error(fmt.Errorf("не удалось отправить access токен: %w", err).Error())
	}
}

func (a *HTTPAuthHandlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		a.logger.Error(fmt.Errorf("отсутствует refresh токен: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	_, err = a.tokenManager.Parse(cookie.Value)
	if err != nil {
		a.logger.Error(fmt.Errorf("неккоректный refresh токен: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	if err = a.tokenService.DeleteRefreshToken(a.ctx, cookie.Value); err != nil {
		a.logger.Error(fmt.Errorf("ошибка удаления refresh токена из базы: %w", err).Error())
		coreErrors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "auth/refresh",
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
}
