package main

import (
	"auth-service/core/enums"
	logs "auth-service/core/logger"
	"auth-service/core/middleware"
	"auth-service/internal/admin/services"
	transport3 "auth-service/internal/admin/transport"
	"auth-service/internal/auth/repositories/interfaces"
	servicesInterfaces "auth-service/internal/auth/services/interfaces"
	"auth-service/internal/auth/transport"
	profileServiceInterfaces "auth-service/internal/designerProfiles/services/interfaces"
	transport2 "auth-service/internal/designerProfiles/transport"
	"context"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "auth-service/docs"

	"github.com/gorilla/mux"
)

type ServerManager struct {
	authHandlers            transport.HTTPAuthHandlers
	designerProfileHandlers transport2.HTTPDesignerProfileHandlers
	adminHandlers           transport3.HTTPAdminHandlers
	middleware              middleware.AuthMiddleware
	logger                  logs.Logger
}

func NewServerManager(tokenService servicesInterfaces.TokenService,
	userService servicesInterfaces.UserService,
	designerProfileService profileServiceInterfaces.DesignerProfileService,
	adminService services.AdminService,
	tokenManager interfaces.TokenManager,
	ctx context.Context,
	logger logs.Logger) ServerManager {
	return ServerManager{
		authHandlers:            transport.NewHTTPAuthHandler(tokenService, tokenManager, userService, ctx, logger),
		designerProfileHandlers: transport2.NewHTTPDesignerProfileHandlers(designerProfileService, ctx, logger),
		adminHandlers:           transport3.NewHTTPAdminHandlers(adminService, ctx, logger),
		middleware:              middleware.NewAuthMiddleware(tokenManager),
		logger:                  logger,
	}
}

func (s *ServerManager) StartServer() {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").
		Handler(httpSwagger.WrapHandler)

	router.
		Path("/auth/register").
		Methods("POST").
		HandlerFunc(
			s.authHandlers.HandleRegister,
		)

	router.
		Path("/auth/login").
		Methods("POST").
		HandlerFunc(
			s.authHandlers.HandleLogin,
		)

	router.
		Path("/auth/refresh").
		Methods("POST").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.authHandlers.HandleRefresh)),
		)

	router.
		Path("/auth/logout").
		Methods("POST").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.authHandlers.HandleLogout)),
		)

	router.
		Path("/designer_profiles/me").
		Methods("POST").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.User,
				http.HandlerFunc(s.designerProfileHandlers.HandleCreateMyDesignerProfile)),
		)

	router.
		Path("/designer_profiles/me").
		Methods("GET").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Designer,
				http.HandlerFunc(s.designerProfileHandlers.HandleGetMyDesignerProfile)),
		)

	router.
		Path("/designer_profiles/me").
		Methods("PATCH").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Designer,
				http.HandlerFunc(s.designerProfileHandlers.HandleUpdateMyDesignerProfile)),
		)

	router.
		Path("/designer_profiles/me").
		Methods("DELETE").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Designer,
				http.HandlerFunc(s.designerProfileHandlers.HandleDeleteMyDesignerProfile)),
		)

	router.
		Path("/role/{user_id}").
		Methods("PATCH").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Admin,
				http.HandlerFunc(s.adminHandlers.HandleChangeRole)),
		)

	router.
		Path("/role/{user_id}").
		Methods("GET").
		HandlerFunc(
			s.middleware.RequireRole(
				enums.Admin,
				http.HandlerFunc(s.adminHandlers.HandleGetRole)),
		)

	err := http.ListenAndServe(":9011", router)
	if err != nil {
		s.logger.Error(fmt.Errorf("ошибка при запуске сервера, %w", err).Error())
		return
	}
}
