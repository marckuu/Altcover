package transport

import (
	coreErrors "auth-service/core/errors"
	logs "auth-service/core/logger"
	globInterfaces "auth-service/internal/auth/repositories/interfaces"
	servicesInterfaces "auth-service/internal/auth/services/interfaces"
	"auth-service/internal/auth/transport/dto"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	if err := a.userService.Register(a.ctx, loginRequest); err != nil {
		a.logger.Error(err.Error())
		coreErrors.SendErrorResponse(w, err, http.StatusInternalServerError)
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

	tokenPair, err := a.userService.Login(a.ctx, loginRequest)
	if err != nil {
		a.logger.Error(err.Error())
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

	accessToken, err := a.userService.Refresh(a.ctx, cookie.Value)

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

	if err = a.userService.Logout(a.ctx, cookie.Value); err != nil {
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
