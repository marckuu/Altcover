package transport

import (
	logs "auth-service/core/logger"
	repositoryInterfaces "auth-service/internal/auth/repositories/interfaces"
	servicesInterfaces "auth-service/internal/auth/services/interfaces"
	"auth-service/internal/auth/transport/dto"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPAuthHandlers_HandleRegister(t *testing.T) {
	userService := servicesInterfaces.NewMockUserService(t)
	tokenManager := repositoryInterfaces.NewMockTokenManager(t)
	tokenService := servicesInterfaces.NewMockTokenService(t)
	logger := logs.NewMockLogger(t)
	ctx := context.Background()

	loginData := dto.LoginRequest{
		Nickname: "User",
		Password: "12345678",
	}

	userService.
		On("Register", ctx, loginData).
		Return(nil)

	body, err := json.Marshal(loginData)
	require.NoError(t, err)

	r := httptest.NewRequest(
		http.MethodPost,
		"/auth/register",
		bytes.NewReader(body),
	)

	w := httptest.NewRecorder()

	a := NewHTTPAuthHandler(
		tokenService,
		tokenManager,
		userService,
		ctx,
		logger,
	)

	a.HandleRegister(w, r)

	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestHTTPAuthHandlers_HandleLogin(t *testing.T) {
	userService := servicesInterfaces.NewMockUserService(t)
	tokenManager := repositoryInterfaces.NewMockTokenManager(t)
	tokenService := servicesInterfaces.NewMockTokenService(t)
	logger := logs.NewMockLogger(t)
	ctx := context.Background()

	loginData := dto.LoginRequest{
		Nickname: "User",
		Password: "12345678",
	}

	userService.
		On("Login", ctx, loginData).
		Return(&repositoryInterfaces.TokenPair{
			AccessToken:  "accessToken",
			RefreshToken: "refreshToken",
		},
			nil)

	body, err := json.Marshal(loginData)
	require.NoError(t, err)

	a := NewHTTPAuthHandler(
		tokenService,
		tokenManager,
		userService,
		ctx,
		logger,
	)

	r := httptest.NewRequest(
		http.MethodPost,
		"/auth/login",
		bytes.NewReader(body),
	)

	w := httptest.NewRecorder()

	a.HandleLogin(w, r)

	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHTTPAuthHandlers_HandleRefresh(t *testing.T) {
	userService := servicesInterfaces.NewMockUserService(t)
	tokenManager := repositoryInterfaces.NewMockTokenManager(t)
	tokenService := servicesInterfaces.NewMockTokenService(t)
	logger := logs.NewMockLogger(t)
	ctx := context.Background()

	cookie := &http.Cookie{
		Name:    "refresh_token",
		Value:   "refreshToken",
		Expires: time.Now().Add(time.Hour),
	}

	userService.
		On("Refresh", ctx, cookie.Value).
		Return("accessToken", nil)

	r := httptest.NewRequest(
		http.MethodPost,
		"/auth/refresh",
		nil,
	)

	w := httptest.NewRecorder()

	r.AddCookie(cookie)

	a := NewHTTPAuthHandler(
		tokenService,
		tokenManager,
		userService,
		ctx,
		logger,
	)

	a.HandleRefresh(w, r)

	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHTTPAuthHandlers_HandleLogout(t *testing.T) {
	userService := servicesInterfaces.NewMockUserService(t)
	tokenManager := repositoryInterfaces.NewMockTokenManager(t)
	tokenService := servicesInterfaces.NewMockTokenService(t)
	logger := logs.NewMockLogger(t)
	ctx := context.Background()

	cookie := &http.Cookie{
		Name:    "refresh_token",
		Value:   "refreshToken",
		Expires: time.Now().Add(time.Hour),
	}

	userService.
		On("Logout", ctx, cookie.Value).
		Return(nil)

	r := httptest.NewRequest(
		http.MethodPost,
		"/auth/logout",
		nil,
	)

	w := httptest.NewRecorder()

	r.AddCookie(cookie)

	a := NewHTTPAuthHandler(
		tokenService,
		tokenManager,
		userService,
		ctx,
		logger,
	)

	a.HandleLogout(w, r)

	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
