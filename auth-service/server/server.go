package server

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"auth-service/repositories"
	"auth-service/services"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerManager struct {
	authHandlers            handlers.HTTPAuthHandlers
	designerProfileHandlers handlers.HTTPDesignerProfileHandlers
	middleware              middleware.AuthMiddleware
}

func NewServerManager(tokenService services.TokenService,
	userService services.UserService,
	designerProfileService services.DesignerProfileService,
	JWTManager repositories.JWTManager,
	ctx context.Context) ServerManager {
	return ServerManager{
		authHandlers:            handlers.NewHTTPAuthHandler(tokenService, JWTManager, userService, ctx),
		designerProfileHandlers: handlers.NewHTTPDesignerProfileHandlers(designerProfileService, ctx),
		middleware:              middleware.NewAuthMiddleware(JWTManager),
	}
}

func (s *ServerManager) StartServer() {
	router := mux.NewRouter()

	router.
		Path("/auth/register").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.authHandlers.HandleRegister)),
		)

	router.
		Path("/auth/login").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.authHandlers.HandleLogin)),
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
		Path("/designer_profiles").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleCreateDesignerProfile)),
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
