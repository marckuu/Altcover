package server

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerManager struct {
	authHandlers            handlers.HTTPAuthHandlers
	designerProfileHandlers handlers.HTTPDesignerProfileHandlers
	middleware              middleware.AuthMiddleware
}

func NewServerManager() ServerManager {
	return ServerManager{
		authHandlers:            handlers.NewHTTPAuthHandler(),
		designerProfileHandlers: handlers.NewHTTPDesignerProfileHandlers(),
		middleware:              middleware.NewAuthMiddleware(),
	}
}

func (s *ServerManager) StartServer() {
	router := mux.NewRouter()

	router.Path("/auth/register").Methods("POST").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.authHandlers.HandleRegister)))
	router.Path("/auth/login").Methods("POST").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.authHandlers.HandleLogin)))
	router.Path("/auth/refresh").Methods("POST").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.authHandlers.HandleRefresh)))
	router.Path("/auth/logout").Methods("POST").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.authHandlers.HandleLogout)))

	router.Path("/designer_profile").Methods("POST").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleCreateDesignerProfile)))
	router.Path("/designer_profile/me").Methods("GET").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleGetCurrentDesignerProfile)))
	router.Path("/designer_profile/me").Methods("PATCH").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleUpdateCurrentDesignerProfile)))
	router.Path("/designer_profile/me").Methods("DELETE").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.designerProfileHandlers.HandleDeleteCurrentDesignerProfile)))

	err := http.ListenAndServe(":9011", router)
	if err != nil {
		println("Ошибка при запуске сервера")
		return
	}
}
