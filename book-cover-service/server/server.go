package server

import (
	"book-cover-service/db/repositories"
	"book-cover-service/handlers"
	"book-cover-service/middleware"
	"book-cover-service/services"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerManager struct {
	coverHandlers     handlers.HTTPCoverHandlers
	coverLikeHandlers handlers.HTTPCoverLikeHandlers
	favoritesHandlers handlers.HTTPFavoritesHandlers
	middleware        middleware.AuthMiddleware
}

func NewServerManager(coverService services.CoverService,
	coverLikeService services.CoverLikeService,
	favoritesService services.FavoritesService,
	JWTManager repositories.JWTManager, ctx context.Context) ServerManager {
	return ServerManager{
		coverHandlers:     handlers.NewCoverHandlers(coverService, ctx),
		coverLikeHandlers: handlers.NewCoverLikeHandlers(coverLikeService, ctx),
		favoritesHandlers: handlers.NewHTTPFavoritesHandlers(favoritesService, ctx),
		middleware:        middleware.NewAuthMiddleware(JWTManager),
	}
}

func (s *ServerManager) StartServer() {
	router := mux.NewRouter()

	router.
		Path("designer/me/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.coverHandlers.HandleGetMyCoversAsDesigner)),
		)

	router.
		Path("designers/{user_id}/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.Auth(
				http.HandlerFunc(s.coverHandlers.HandleGetCoversByDesignerUserID),
			),
		)

	router.
		Path("/covers/{cover_id}").
		Methods("PATCH").
		HandlerFunc(
			s.middleware.Auth(
				http.HandlerFunc(s.coverHandlers.HandleUpdateCover),
			),
		)

	router.
		Path("/covers/{cover_id}").
		Methods("GET").
		HandlerFunc(s.coverHandlers.HandleGetCoverByID)

	// Дописать внутри запрос на получение данных профиля
	router.Path("/covers").Methods("POST").HandlerFunc(s.middleware.Auth(http.HandlerFunc(s.coverHandlers.HandleAddCover)))

	router.
		Path("feeds/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(s.coverHandlers.HandleGetFeedCovers)

	router.
		Path("books/{book_id}/covers").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(s.coverHandlers.HandleGetCoversByBook)

	router.
		Path("/covers/{cover_id}/likes").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.coverLikeHandlers.HandleSetLike)),
		)

	router.
		Path("/covers{cover_id}/favorites").
		Methods("POST").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.favoritesHandlers.HandleAddCoverToFavorites)),
		)

	router.
		Path("/covers/me/favorites").
		Methods("GET").
		Queries("offset", "{offset}", "limit", "{limit}").
		HandlerFunc(
			s.middleware.Auth(http.HandlerFunc(s.favoritesHandlers.HandleGetMyFavoriteCovers)),
		)

	err := http.ListenAndServe(":9011", router)
	if err != nil {
		println("Ошибка при запуске сервера")
		return
	}
}
