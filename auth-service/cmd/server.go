package main

import (
	logs "auth-service/core/logger"
	"auth-service/core/middleware"
	"auth-service/internal/auth/repositories/interfaces"
	servicesInterfaces "auth-service/internal/auth/services/interfaces"
	"auth-service/internal/auth/transport"
	profileServiceInterfaces "auth-service/internal/designerProfiles/services/interfaces"
	transport2 "auth-service/internal/designerProfiles/transport"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerManager struct {
	authHandlers            transport.HTTPAuthHandlers
	designerProfileHandlers transport2.HTTPDesignerProfileHandlers
	middleware              middleware.AuthMiddleware
}

func NewServerManager(tokenService servicesInterfaces.TokenService,
	userService servicesInterfaces.UserService,
	designerProfileService profileServiceInterfaces.DesignerProfileService,
	tokenManager interfaces.TokenManager,
	ctx context.Context,
	logger logs.Logger) ServerManager {
	return ServerManager{
		authHandlers:            transport.NewHTTPAuthHandler(tokenService, tokenManager, userService, ctx, logger),
		designerProfileHandlers: transport2.NewHTTPDesignerProfileHandlers(designerProfileService, ctx, logger),
		middleware:              middleware.NewAuthMiddleware(tokenManager),
	}
}

func (s *ServerManager) StartServer() {
	router := mux.NewRouter()

	router.
		Path("/auth/register").
		Methods("POST").
		HandlerFunc(
			http.HandlerFunc(s.authHandlers.HandleRegister),
		)

	router.
		Path("/auth/login").
		Methods("POST").
		HandlerFunc(
			http.HandlerFunc(s.authHandlers.HandleLogin),
		)

	router.
		Path("/auth/refresh").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.authHandlers.HandleRefresh)),
		)

	router.
		Path("/auth/logout").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.authHandlers.HandleLogout)),
		)

	router.
		Path("/designer_profiles/me").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleCreateMyDesignerProfile)),
		)

	router.
		Path("/designer_profiles/me").
		Methods("GET").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleGetMyDesignerProfile)),
		)

	router.
		Path("/designer_profiles/me").
		Methods("PATCH").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleUpdateMyDesignerProfile)),
		)

	router.
		Path("/designer_profiles/me").
		Methods("DELETE").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleDeleteMyDesignerProfile)),
		)

	err := http.ListenAndServe(":9011", router)
	if err != nil {
		println("Ошибка при запуске сервера")
		return
	}
}
